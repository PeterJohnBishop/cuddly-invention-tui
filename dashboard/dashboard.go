package dashboard

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// model
type DashboardModel struct {
	cursor             int
	selected           map[int]struct{}
	status             bool
	err                error
	cpuPercent         []float64
	totalVM            uint64
	freeVM             uint64
	usedPercentVM      float64
	totalSwapMem       uint64
	freeSwapMem        uint64
	usedPercentSwapMem float64
	totalDisk          uint64
	freeDisk           uint64
	usedPercentDisk    float64
	dockerClient       *client.Client
	containers         []container.Summary
	containerMetrics   map[string]ContainerStats
}

type ContainerStats struct {
	CPU    float64
	Memory float64
}

func InitialDashboardModel() DashboardModel {

	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return DashboardModel{
		selected:           make(map[int]struct{}),
		status:             false,
		err:                nil,
		cpuPercent:         nil,
		totalVM:            0,
		freeVM:             0,
		usedPercentVM:      0,
		totalSwapMem:       0,
		freeSwapMem:        0,
		usedPercentSwapMem: 0,
		totalDisk:          0,
		freeDisk:           0,
		usedPercentDisk:    0,
		dockerClient:       cli,
		containers:         []container.Summary{},
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.tickCPU(),
		m.tickMemory(),
		m.tickSwap(),
		m.tickDisk(),
		m.tickDocker(),
	)
}

// update method
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case errMsg:
		m.err = msg.Err
		return m, nil

	case cpuUsageMsg:
		m.cpuPercent = msg.Percent
		return m, m.tickCPU()

	case virtualMemoryMsg:
		m.totalVM = msg.Total
		m.freeVM = msg.Free
		m.usedPercentVM = msg.UsedPercent
		return m, m.tickMemory()

	case swapMemoryMsg:
		m.totalSwapMem = msg.Total
		m.freeSwapMem = msg.Free
		m.usedPercentSwapMem = msg.UsedPercent
		return m, m.tickSwap()

	case diskMsg:
		m.totalDisk = msg.Total
		m.freeDisk = msg.Free
		m.usedPercentDisk = msg.UsedPercent
		return m, m.tickDisk()

	case dockerMsg:
		m.containers = msg.Containers
		return m, m.tickDocker()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.containers)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

// view method
func (m DashboardModel) View() string {
	var b strings.Builder

	b.WriteString("\n\n[ SYSTEM DASHBOARD ]\n\n")

	// (CPU, RAM, Swap, Disk...)
	cpuVal := 0.0
	if len(m.cpuPercent) > 0 {
		cpuVal = m.cpuPercent[0]
	}
	b.WriteString(fmt.Sprintf("CPU  : %-20s %5.1f%% \n", makeBar(cpuVal), cpuVal))
	b.WriteString(fmt.Sprintf("RAM  : %-20s %5.1f%% \n", makeBar(m.usedPercentVM), m.usedPercentVM))
	b.WriteString(fmt.Sprintf("SWAP : %-20s %5.1f%% \n", makeBar(m.usedPercentSwapMem), m.usedPercentSwapMem))
	b.WriteString(fmt.Sprintf("DISK : %-20s %5.1f%% \n", makeBar(m.usedPercentDisk), m.usedPercentDisk))

	b.WriteString("\n\n[ DOCKER CONTAINERS ]\n\n")
	b.WriteString(fmt.Sprintf("      %-12s %-20s %-10s\n", "ID", "IMAGE", "STATUS"))

	if len(m.containers) == 0 {
		b.WriteString("\n[ No containers running. ]\n")
	}

	for i, c := range m.containers {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		shortID := c.ID[:10]
		image := c.Image
		if len(image) > 18 {
			image = image[:15] + "..."
		}

		status := c.State // "running", "exited", etc.

		b.WriteString(fmt.Sprintf("%s [%s] %-12s %-20s %-10s\n", cursor, checked, shortID, image, status))
		if checked == "x" {
			stats := m.containerMetrics[c.ID]
			statsStr := fmt.Sprintf("CPU: %.1f%% | MEM: %.1f%%\n", stats.CPU, stats.Memory)
			b.WriteString(fmt.Sprintf("      - %s", statsStr))
		}
	}

	// err handling
	if m.err != nil {
		b.WriteString("\nStatus: " + m.err.Error() + "\n")
		b.WriteString(" (Press 'esc' to dismiss)\n")
	} else {
		b.WriteString("\nStatus: Running\n")
		b.WriteString("[up/k] [down/j] to move cursor | [enter/space] to select\n")
		b.WriteString("[ctrl+c] or [q] to quit")
	}

	return b.String()
}

// simple ASCII progress bar
func makeBar(percent float64) string {
	width := 20
	fullBlocks := int(percent / 100 * float64(width))
	if fullBlocks > width {
		fullBlocks = width
	}

	return "[" + strings.Repeat("█", fullBlocks) + strings.Repeat(" ", width-fullBlocks) + "]"
}
