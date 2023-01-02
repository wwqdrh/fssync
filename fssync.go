package main

import (
	"github.com/wwqdrh/fssync/cli"
)

// var mode = flag.String("mode", "", "client or server")

func main() {
	cli.StartCli()

	// flag.Parse()
	// if *mode == "" {
	// 	flag.Usage()
	// 	return
	// }

	// switch *mode {
	// case "client":
	// 	if err := tea.NewProgram(client.NewClientView()).Start(); err != nil {
	// 		fmt.Printf("could not start program: %s\n", err)
	// 		os.Exit(1)
	// 	}
	// case "server":
	// 	server.StartCli()
	// }
}
