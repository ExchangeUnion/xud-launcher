package cmd

import (
	"github.com/spf13/cobra"
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

	},
}
