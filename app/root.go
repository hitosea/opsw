package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootApp = &cobra.Command{
	Use:   "opsw",
	Short: "OPSW is a tool for managing your servers",
}

func Execute() {
	rootApp.CompletionOptions.DisableDefaultCmd = true
	if err := rootApp.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
