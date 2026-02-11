package dashboard

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// model
type DashboardModel struct {
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
}

func InitialDashboardModel() DashboardModel {
	return DashboardModel{
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
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.tickCPU(),
		m.tickMemory(),
		m.tickSwap(),
		m.tickDisk(),
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

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// view method
func (m DashboardModel) View() string {
	var b strings.Builder

	b.WriteString("╭─────────── SYSTEM DASHBOARD ───────────╮\n")

	// (CPU, RAM, Swap, Disk...)
	cpuVal := 0.0
	if len(m.cpuPercent) > 0 {
		cpuVal = m.cpuPercent[0]
	}
	b.WriteString(fmt.Sprintf("│ CPU  : %-20s %5.1f%% │\n", makeBar(cpuVal), cpuVal))
	b.WriteString(fmt.Sprintf("│ RAM  : %-20s %5.1f%% │\n", makeBar(m.usedPercentVM), m.usedPercentVM))
	b.WriteString(fmt.Sprintf("│ SWAP : %-20s %5.1f%% │\n", makeBar(m.usedPercentSwapMem), m.usedPercentSwapMem))
	b.WriteString(fmt.Sprintf("│ DISK : %-20s %5.1f%% │\n", makeBar(m.usedPercentDisk), m.usedPercentDisk))

	b.WriteString("╰────────────────────────────────────────╯\n")

	// err handling
	if m.err != nil {
		b.WriteString("\n ❌ ERROR: " + m.err.Error() + "\n")
		b.WriteString(" (Press 'esc' to dismiss)")
	} else {
		b.WriteString("\n [q] Quit | Status: Running")
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
