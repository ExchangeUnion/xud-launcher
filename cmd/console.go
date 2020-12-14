package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(consoleCmd)
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open your native console with some useful commands for xud-docker",
	Run: func(cmd *cobra.Command, args []string) {
		exec.Command("/bin/bash")
	},
}
