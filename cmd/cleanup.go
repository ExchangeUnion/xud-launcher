package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func isValidNetwork(network string) bool {
	network = strings.TrimSpace(network)
	network = strings.ToLower(network)
	for _, item := range validNetworks {
		if item == network {
			return true
		}
	}
	return false
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup the XUD environment",
	//ValidArgs: []string{
	//	"network",
	//},
	Args: func(cmd *cobra.Command, args []string) error {
		//err := cobra.OnlyValidArgs(cmd, args)
		//if err != nil {
		//	return err
		//}
		if len(args) < 1 {
			return errors.New("network required")
		}
		if !isValidNetwork(args[0]) {
			return fmt.Errorf("invalid network: %s", args[0])
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]

		// get running containers
		containers := getRunningContainers(network)
		if len(containers) > 0 {
			fmt.Println("Stopping containers...")
			for _, c := range containers {
				stopContainer(c)
			}
		}

		containers = getContainers(network)
		if len(containers) > 0 {
			fmt.Println("Removing containers...")
			for _, c := range containers {
				removeContainer(c)
			}
		}

		networks := getNetworks(network)
		if len(containers) > 0 {
			fmt.Println("Removing networks...")
			for _, n := range networks {
				removeNetwork(n)
			}
		}

		home, err := homedir.Dir()
		if err != nil {
			logger.Fatal(err)
		}
		networkDir := filepath.Join(home, ".xud-docker", network)
		fmt.Println(networkDir)

		if _, err := os.Stat(networkDir); !os.IsNotExist(err) {
			fmt.Printf("Do you want to remove all %s data (%s)? [y/N] ", network, networkDir)
			var reply string
			_, err := fmt.Scanln(&reply)
			if err != nil {
				logger.Fatal(err)
			}
			reply = strings.ToLower(reply)
			if reply == "y" || reply == "yes" {
				fmt.Println("Removing data...")
				removeDir(networkDir)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
}

func getRunningContainers(network string) []string {
	var result []string
	filter := fmt.Sprintf("name=%s_", network)
	out, err := exec.Command("docker", "ps", "--filter", filter, "--format", "{{.ID}}").Output()
	if err != nil {
		logger.Fatal(err)
	}
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		result = append(result, s.Text())
	}
	return result
}

func getContainers(network string) []string {
	var result []string
	filter := fmt.Sprintf("name=%s_", network)
	out, err := exec.Command("docker", "ps", "--filter", filter, "--format", "{{.ID}}", "-a").Output()
	if err != nil {
		logger.Fatal(err)
	}
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		result = append(result, s.Text())
	}
	return result
}

func getNetworks(network string) []string {
	var result []string
	filter := fmt.Sprintf("name=%s_", network)
	out, err := exec.Command("docker", "network", "ls", "--filter", filter, "--format", "{{.ID}}").Output()
	if err != nil {
		logger.Fatal(err)
	}
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		result = append(result, s.Text())
	}
	return result
}

func stopContainer(id string) {
	fmt.Println(id)
	err := exec.Command("docker", "stop", id).Run()
	if err != nil {
		logger.Fatal(err)
	}
}

func removeContainer(id string) {
	fmt.Println(id)
	err := exec.Command("docker", "rm", id).Run()
	if err != nil {
		logger.Fatal(err)
	}
}

func removeNetwork(id string) {
	fmt.Println(id)
	err := exec.Command("docker", "network", "rm", id).Run()
	if err != nil {
		logger.Fatal(err)
	}
}

func removeDir(path string) {
	fmt.Println(path)
	err := exec.Command("sudo", "rm", "-rf", path).Run()
	if err != nil {
		logger.Fatal(err)
	}
}
