package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print the basic information of xud-launcher and your xud-docker setups",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO implement the info command
		// return
		// * xud-launcher version
		// * xud-docker network dir location
		// * xud-docker installed or not
		fmt.Println("To be implemented!")
	},
}
