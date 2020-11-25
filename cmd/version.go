package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of xud-launcher",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.0-commit")
	},
}
