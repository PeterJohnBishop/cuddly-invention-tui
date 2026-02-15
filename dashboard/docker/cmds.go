package docker

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/moby/moby/client"
)

func (m DockerModel) tickDocker() tea.Cmd {
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
