package dashboard

import (
	"cuddly-invention-tui/dashboard/docker"
	"cuddly-invention-tui/dashboard/kubernetes"
	"cuddly-invention-tui/dashboard/system"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	width      int
	height     int
	system     system.SystemModel
	docker     docker.DockerModel
	kubernetes kubernetes.KubernetesModel
}

func InitialDashboardModel() DashboardModel {
	return DashboardModel{
		system:     system.InitialSystemModel(),
		docker:     docker.InitialDockerModel(),
		kubernetes: kubernetes.InitalKubernetsModel(),
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.system.Init(),
		m.docker.Init(),
		m.kubernetes.Init(),
	)
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	m.system, cmd = m.system.Update(msg)
	cmds = append(cmds, cmd)

	m.docker, cmd = m.docker.Update(msg)
	cmds = append(cmds, cmd)

	m.kubernetes, cmd = m.kubernetes.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m DashboardModel) View() string {
	rowWidth := m.width - 2

	style := lipgloss.NewStyle().
		Width(rowWidth).
		Border(lipgloss.NormalBorder())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		style.Render(m.system.View()),
		style.Render(m.docker.View()),
		style.Render(m.kubernetes.View()),
	)
}
