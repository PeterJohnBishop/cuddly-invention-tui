package docker

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type DockerModel struct {
	cursor           int
	selected         map[int]struct{}
	status           bool
	err              error
	dockerClient     *client.Client
	dockerMem        int64
	containers       []container.Summary
	containerMetrics map[string]ContainerStats
}

func InitialDockerModel() DockerModel {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	opt := client.InfoOptions{}
	info, err := cli.Info(ctx, opt)
	if err != nil {
		panic(err)
	}

	return DockerModel{
		selected:     make(map[int]struct{}),
		status:       false,
		err:          nil,
		dockerClient: cli,
		dockerMem:    info.Info.MemTotal / 1024 / 1024 / 1024,
		containers:   []container.Summary{},
	}
}

func (m DockerModel) Init() tea.Cmd {
	return tea.Batch(
		m.tickDocker(),
	)
}

func (m DockerModel) Update(msg tea.Msg) (DockerModel, tea.Cmd) {
	switch msg := msg.(type) {

	case errMsg:
		m.err = msg.Err
		return m, nil

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

func (m DockerModel) View() string {
	var b strings.Builder
	b.WriteString("\n\n[ DOCKER ]\n\n")
	b.WriteString(fmt.Sprintf("System memory assigned to Docker: %d GB\n\n", m.dockerMem))
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

	if m.err != nil {
		b.WriteString("\nStatus: " + m.err.Error() + "\n")
		b.WriteString(" (Press 'esc' to dismiss)\n")
	}

	b.WriteString("[up/k] [down/j] to move cursor | [enter/space] to select\n")
	b.WriteString("[ctrl+c] or [q] to quit")

	return b.String()
}
