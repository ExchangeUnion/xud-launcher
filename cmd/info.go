package cmd

import (
	"fmt"
	"github.com/ExchangeUnion/xud-launcher/build"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print the basic information of xud-launcher and your xud-docker setups",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Launcher Version: %s\n", GetVersion())
		fmt.Printf("Built At: %s\n", GetBuildTime())
		fmt.Printf("XUD-Docker Version: 20.12.18-01\n")
		fmt.Printf("Network: %s\n", network)
		fmt.Printf("Network Directory: %s\n", networkDir)
	},
}

func GetBuildTime() string {
	i, err := strconv.ParseInt(build.Timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	return fmt.Sprint(tm)
}
