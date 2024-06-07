package main

import (
	"fmt"
	"os"
	"path"

	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/fssync/server"
	"github.com/wwqdrh/gokit/clitool"
)

func initSpecPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	specPath := path.Join(home, ".fssync")
	if _, err := os.Stat(specPath); os.IsNotExist(err) {
		if err := os.MkdirAll(specPath, os.ModePerm); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	return specPath
}

func main() {
	client.RootSpecPath = initSpecPath()
	cmd := clitool.Command{}
	cmd.Add(client.Command())
	cmd.Add(server.Command())
	cmd.Run()
}
