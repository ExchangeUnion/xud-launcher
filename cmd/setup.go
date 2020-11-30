package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	restore   bool
	backupDir string
)

func init() {
	setupCmd.PersistentFlags().String("wallet-password", "", "XUD master wallet password")
	setupCmd.PersistentFlags().StringVar(&backupDir, "backup-dir", "", "XUD backup location")
	setupCmd.PersistentFlags().BoolVar(&restore, "restore", true, "Restore wallets")

	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Easy setup flow to bring up your XUD environment",
	Run: func(cmd *cobra.Command, args []string) {
		// start proxy(api) service first
		upService("proxy")

		// start L2 services like lndbtc, lndltc and connext
		upService("lndbtc")
		upService("lndltc")
		upService("connext")

		// waiting for lnds to be synced

		// start XUD service
		if !hasBackup() {
			var reason string

			if backupDir == "" {
				for {
					backupDir := input("Enter path to backup location")

					if checkDir(backupDir, &reason) {
						break
					} else {
						fmt.Printf("invalid backup location (%s): %s\n", reason, backupDir)
					}
				}
			} else {
				if !checkDir(backupDir, &reason) {
					fmt.Printf("invalid backup location (%s): %s\n", reason, backupDir)
					os.Exit(1)
				}
			}

			setBackup(backupDir)
		}
		upService("xud")

		// post actions like: create/restore wallets
		if restore {
			serviceExec("xud", "xucli", "restore")
		} else {
			serviceExec("xud", "xucli", "create")
		}

		// start optional services like: arby, boltz, webui
		upService("arby")
		upService("boltz")
		upService("webui")
	},
}

func checkDir(path string, reason *string) bool {
	return true
}

func input(prompt string) string {
	fmt.Printf("%s: ", prompt)
	var reply string
	_, err := fmt.Scanln(&reply)
	if err != nil {
		logger.Fatal(err)
	}
	return reply
}

func yesNo(question string) string {
	fmt.Printf("%s ? [Y/n] ", question)
	var reply string
	_, err := fmt.Scanln(&reply)
	if err != nil {
		logger.Fatal(err)
	}
	reply = strings.ToLower(reply)
	if reply == "y" || reply == "yes" || reply == "" {
		return "yes"
	}
	return "no"
}

func noYes(question string) string {
	fmt.Printf("%s ? [y/N] ", question)
	var reply string
	_, err := fmt.Scanln(&reply)
	if err != nil {
		logger.Fatal(err)
	}
	reply = strings.ToLower(reply)
	if reply == "n" || reply == "no" || reply == "" {
		return "no"
	}
	return "yes"
}

func upService(name string) {
	upCmd.SetArgs([]string{"-d", name})
	if err := upCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func serviceExec(service string, args ...string) {

}

func hasBackup() bool {
	return false
}

func setBackup(path string) {

}
