package dashboard

import tea "github.com/charmbracelet/bubbletea"

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

func initialDashboardModel() DashboardModel {
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
	return nil
}

// update method
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cpuUsageMsg:
		m.cpuPercent = msg.Percent
	case virtualMemoryMsg:
		m.totalVM = msg.Total
		m.freeVM = msg.Free
		m.usedPercentVM = msg.UsedPercent
	case swapMemoryMsg:
		m.totalSwapMem = msg.Total
		m.freeSwapMem = msg.Free
		m.usedPercentSwapMem = msg.UsedPercent
	case diskMsg:
		m.totalDisk = msg.Total
		m.freeDisk = msg.Free
		m.usedPercentDisk = msg.UsedPercent
	}
	return m, nil
}

// view method
func (m DashboardModel) View() string {
	s := ""
	return s
}
