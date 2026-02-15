package kubernetes

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type KubernetesModel struct {
	status bool
}

func InitalKubernetsModel() KubernetesModel {
	return KubernetesModel{
		status: false,
	}
}
func (m KubernetesModel) Init() tea.Cmd { return nil }
func (m KubernetesModel) Update(msg tea.Msg) (KubernetesModel, tea.Cmd) {
	return m, nil
}
func (m KubernetesModel) View() string {
	var status string
	if m.status {
		status = "up"
	} else {
		status = "down"
	}
	return fmt.Sprintf("kubernetesModel \nStatus: %s", status)
}
