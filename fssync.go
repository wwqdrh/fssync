package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wwqdrh/fssync/cli/client"
	"github.com/wwqdrh/fssync/cli/server"
)

var mode = flag.String("mode", "", "client or server")

func main() {
	flag.Parse()
	if *mode == "" {
		flag.Usage()
		return
	}

	switch *mode {
	case "client":
		if err := tea.NewProgram(client.NewClientView()).Start(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	case "server":
		if err := tea.NewProgram(server.NewServerView()).Start(); err != nil {
			fmt.Printf("could not start program: %s\n", err)
			os.Exit(1)
		}
	}
}
