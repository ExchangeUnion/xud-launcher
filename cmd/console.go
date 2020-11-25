package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(consoleCmd)
}

var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Start the xud-ctl console",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting...")
	},
}
