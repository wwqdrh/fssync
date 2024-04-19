package main

import (
	"github.com/wwqdrh/fssync/client"
	"github.com/wwqdrh/fssync/server"
	"github.com/wwqdrh/gokit/clitool"
)

func main() {
	cmd := clitool.Command{}
	cmd.Add(client.Command())
	cmd.Add(server.Command())
	cmd.Run()
}
