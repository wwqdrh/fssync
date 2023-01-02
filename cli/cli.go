package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wwqdrh/fssync/cli/client"
	"github.com/wwqdrh/fssync/cli/server"
	"github.com/wwqdrh/gokit/clitool"
)

var Opt = struct {
	Name string
}{}

func StartCli() {
	cmd := clitool.Command{
		Cmd: &cobra.Command{
			Use:   "connect",
			Short: "Create a network tunnel to docker swarm cluster",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("do run", Opt.Name)
				return nil
			},
			Example: "ktctl connect [command options]",
		},
		Persistent: []clitool.OptionConfig{
			{
				Target:       "Name",
				Name:         "name",
				DefaultValue: "127.0.0.1:18080",
				Description:  "docker swarm集群中的shadow服务地址",
			},
		},
		Values: &Opt,
	}
	cmd.Add(client.Command())
	cmd.Add(server.Command())

	cmd.Run()
}
