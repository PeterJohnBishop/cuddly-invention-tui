package dashboard

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moby/moby/client"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func (m DashboardModel) tickCPU() tea.Cmd {
	return tea.Tick(time.Millisecond*800, func(t time.Time) tea.Msg {
		cpuPercent, _ := cpu.Percent(0, false)
		return cpuUsageMsg{Percent: cpuPercent}
	})
}

func (m DashboardModel) tickMemory() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		v, _ := mem.VirtualMemory()
		return virtualMemoryMsg{
			Total:       v.Total / 1024 / 1024 / 1024,
			Free:        v.Free / 1024 / 1024 / 1024,
			UsedPercent: v.UsedPercent,
		}
	})
}

func (m DashboardModel) tickSwap() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		s, err := mem.SwapMemory()
		if err != nil {
			return errMsg{err}
		}
		return swapMemoryMsg{
			Total:       s.Total / 1024 / 1024 / 1024, // GB
			Free:        s.Free / 1024 / 1024 / 1024,  // GB
			UsedPercent: s.UsedPercent,
		}
	})
}

func (m DashboardModel) tickDisk() tea.Cmd {
	return tea.Tick(time.Second*30, func(t time.Time) tea.Msg {
		d, err := disk.Usage("/")
		if err != nil {
			return errMsg{err}
		}
		return diskMsg{
			Total:       d.Total / 1024 / 1024 / 1024, // GB
			Free:        d.Free / 1024 / 1024 / 1024,  // GB
			UsedPercent: d.UsedPercent,
		}
	})
}

func (m DashboardModel) tickDocker() tea.Cmd {
	return tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
		if m.dockerClient == nil {
			return errMsg{Err: fmt.Errorf("Docker client not initialized")}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		containers, err := m.dockerClient.ContainerList(ctx, client.ContainerListOptions{All: true})
		if err != nil {
			return errMsg{Err: err}
		}

		return dockerMsg{Containers: containers.Items}
	})
}
