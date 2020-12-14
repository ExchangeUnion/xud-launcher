package cmd

import (
	"fmt"
	"github.com/ExchangeUnion/xud-launcher/build"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the launcher version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetVersion())
	},
}

func GetVersion() string {
	var version = fmt.Sprintf("%s-%s", build.Version, build.GitCommit[:7])
	return version
}
