package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(pullCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "A docker-compose up wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Chdir(networkDir)
		if err != nil {
			logger.Fatal(err)
		}
		args = append([]string{"up"}, args...)
		c := exec.Command("docker-compose", args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "A docker-compose down wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Chdir(networkDir)
		if err != nil {
			logger.Fatal(err)
		}
		args = append([]string{"down"}, args...)
		c := exec.Command("docker-compose", args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A docker-compose start wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Chdir(networkDir)
		if err != nil {
			logger.Fatal(err)
		}
		args = append([]string{"start"}, args...)
		c := exec.Command("docker-compose", args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A docker-compose stop wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Chdir(networkDir)
		if err != nil {
			logger.Fatal(err)
		}
		args = append([]string{"stop"}, args...)
		c := exec.Command("docker-compose", args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "A docker-compose restart wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Chdir(networkDir)
		if err != nil {
			logger.Fatal(err)
		}
		args = append([]string{"restart"}, args...)
		c := exec.Command("docker-compose", args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "A docker-compose logs wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Chdir(networkDir)
		if err != nil {
			logger.Fatal(err)
		}
		args = append([]string{"logs"}, args...)
		c := exec.Command("docker-compose", args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A docker-compose exec wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Chdir(networkDir)
		if err != nil {
			logger.Fatal(err)
		}
		args = append([]string{"exec"}, args...)
		c := exec.Command("docker-compose", args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Fatal(err)
		}
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A docker-compose pull wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Chdir(networkDir)
		if err != nil {
			logger.Fatal(err)
		}
		args = append([]string{"pull"}, args...)
		c := exec.Command("docker-compose", args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err = c.Run()
		if err != nil {
			logger.Fatal(err)
		}
	},
}
