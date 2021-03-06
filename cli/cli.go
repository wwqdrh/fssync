package cli

import (
	"github.com/spf13/cobra"
	"github.com/wwqdrh/logger"
)

var rootCmd = &cobra.Command{
	Use:               "文件同步工具",
	Short:             "文件同步工具",
	SilenceUsage:      true,
	Long:              `文件同步工具, 提供客户端与服务端`,
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(ServerCmd)
	rootCmd.AddCommand(ClientCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if err := rootCmd.Help(); err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
	}
}
