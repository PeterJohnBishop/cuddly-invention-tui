package main

import (
	"cuddly-invention-tui/dashboard"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Println("Brewing Bubbletea")
	p := tea.NewProgram(dashboard.InitialDashboardModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
