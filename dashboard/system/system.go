package system

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SystemModel struct {
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
	progress           progress.Model
}

func InitialSystemModel() SystemModel {
	return SystemModel{
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
		progress: progress.New(
			progress.WithGradient("#FF00D4", "#01FFEF"),
			progress.WithoutPercentage(),
		),
	}
}

func (m SystemModel) Init() tea.Cmd {
	return tea.Batch(
		m.tickCPU(),
		m.tickMemory(),
		m.tickSwap(),
		m.tickDisk(),
	)
}

func (m SystemModel) Update(msg tea.Msg) (SystemModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 20
		return m, nil

	case progress.FrameMsg:
		newProgressModel, cmd := m.progress.Update(msg)
		m.progress = newProgressModel.(progress.Model)
		return m, cmd

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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m SystemModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).MarginBottom(1)
	labelStyle := lipgloss.NewStyle().Width(6).Bold(true)
	percentStyle := lipgloss.NewStyle().Width(8).Align(lipgloss.Right)

	b.WriteString(headerStyle.Render("[ SYSTEM ]") + "\n\n")

	cpuVal := 0.0
	if len(m.cpuPercent) > 0 {
		cpuVal = m.cpuPercent[0]
	}

	renderRow := func(label string, percent float64) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			labelStyle.Render(label),
			m.progress.ViewAs(percent/100),
			percentStyle.Render(fmt.Sprintf("%5.1f%%", percent)),
		) + "\n"
	}

	b.WriteString(renderRow("CPU", cpuVal))
	b.WriteString(renderRow("RAM", m.usedPercentVM))
	b.WriteString(renderRow("SWAP", m.usedPercentSwapMem))
	b.WriteString(renderRow("DISK", m.usedPercentDisk))

	return b.String()
}
