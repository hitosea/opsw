package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"opsw/vars"
)

func init() {
	rootCommand.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本号",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommitSHA: %s\n", vars.Version, vars.CommitSHA)
	},
}
