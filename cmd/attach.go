package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(attachCmd)
}

var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach this launcher to xud-docker proxy(api) container",
	Long: `The xud-docker proxy(api) container will delegate all container operations to an attached xud-launcher. So 
we don't need to map any docker.sock file into proxy container and it will make the proxy robust across different 
platforms. The proxy will communicate with the launcher through a WebSocket connection by reusing the HTTPS/HTTP port 
exposed to the host`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO implement the attach command
		fmt.Println("To be implemented!")
	},
}
