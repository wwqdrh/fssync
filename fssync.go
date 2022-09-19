package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwqdrh/fssync/cli/tui"
)

func main() {
	if err := tea.NewProgram(tui.NewClientView()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
