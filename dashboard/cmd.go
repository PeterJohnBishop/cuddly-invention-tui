package dashboard

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func checkCPUUsage() tea.Msg {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return errMsg{err}
	}
	return cpuUsageMsg{cpuPercent}
}

func checkVirtualMemory() tea.Msg {
	v, err := mem.VirtualMemory()
	if err != nil {
		return errMsg{err}
	}
	return virtualMemoryMsg{
		Total:       v.Total / 1024 / 1024 / 1024, // GB
		Free:        v.Free / 1024 / 1024 / 1024,  // GB
		UsedPercent: v.UsedPercent,                // %
	}
}

func checkSwapMemory() tea.Msg {
	s, err := mem.SwapMemory()
	if err != nil {
		return errMsg{err}
	}
	return swapMemoryMsg{
		Total:       s.Total,
		Free:        s.Free,
		UsedPercent: s.UsedPercent,
	}
}

func checkDiskUsage() tea.Msg {
	d, err := disk.Usage("/")
	if err != nil {
		return errMsg{err}
	}
	return diskMsg{
		Total:       d.Total,
		Free:        d.Free,
		UsedPercent: d.UsedPercent,
	}
}
