package cmd

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

var (
	logger     = logrus.New()
	network    string
	homeDir    string
	networkDir string
	dataDir    string
	logsDir    string

	rootCmd = &cobra.Command{
		Use:   "xud-launcher",
		Short: "XUD environment launcher",
	}
)

func getHomeDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	switch runtime.GOOS {
	case "linux":
		return filepath.Join(homeDir, ".xud-docker")
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "XudDocker")
	case "windows":
		return filepath.Join(homeDir, "AppData", "Local", "XudDocker")
	default:
		panic(errors.New("unsupported platform: " + runtime.GOOS))
	}
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModeDir|0700); err != nil {
			panic(err)
		}
		logger.Debugf("Created folder: " + path)
	}
}

func init() {
	logger.SetLevel(logrus.DebugLevel)

	homeDir = getHomeDir()
	ensureDir(homeDir)

	cobra.OnInitialize(initConfig)

	logger.Info("Configuring global flags")

	rootCmd.PersistentFlags().StringVarP(&network, "network", "n", "simnet", "specify XUD network")
	rootCmd.PersistentFlags().String("simnet-dir", filepath.Join(homeDir, "simnet"), "Simnet environment folder")
	rootCmd.PersistentFlags().String("testnet-dir", filepath.Join(homeDir, "testnet"), "Testnet environment folder")
	rootCmd.PersistentFlags().String("mainnet-dir", filepath.Join(homeDir, "mainnet"), "Mainnet environment folder")

	if err := viper.BindPFlag("simnet-dir", rootCmd.PersistentFlags().Lookup("simnet-dir")); err != nil {
		logger.Fatal(err)
	}
	if err := viper.BindPFlag("testnet-dir", rootCmd.PersistentFlags().Lookup("testnet-dir")); err != nil {
		logger.Fatal(err)
	}
	if err := viper.BindPFlag("mainnet-dir", rootCmd.PersistentFlags().Lookup("mainnet-dir")); err != nil {
		logger.Fatal(err)
	}
}

func initConfig() {
	generalConf := filepath.Join(homeDir, "xud-docker.conf")
	viper.SetConfigFile(generalConf)
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	logger.Infof("Loading general config file: %s", generalConf)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Info(err)
	}

	switch network {
	case "simnet":
		networkDir = viper.GetString("simnet-dir")
	case "testnet":
		networkDir = viper.GetString("testnet-dir")
	case "mainnet":
		networkDir = viper.GetString("mainnet-dir")
	default:
		panic(errors.New("invalid network: " + network))
	}

	ensureDir(networkDir)

	dataDir := filepath.Join(networkDir, "data")
	ensureDir(dataDir)
	logsDir := filepath.Join(networkDir, "logs")
	ensureDir(logsDir)

	logger.Debugf("homeDir=%s", homeDir)
	logger.Debugf("networkDir=%s", networkDir)
	logger.Debugf("dataDir=%s", dataDir)
	logger.Debugf("logsDir=%s", logsDir)

	networkConf := filepath.Join(networkDir, fmt.Sprintf("%s.conf", network))
	viper.SetConfigFile(networkConf)
	viper.SetConfigType("toml")

	logger.Infof("Loading network config file: %s", networkConf)
	err = viper.MergeInConfig()
	if err != nil {
		logger.Info(err)
	}
}
