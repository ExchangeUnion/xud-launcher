package cmd

import (
	"errors"
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
	simnetDir  string
	testnetDir string
	mainnetDir string
	networkDir string
	dataDir string
	logsDir string



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
		if err := os.Mkdir(path, os.ModeDir); err != nil {
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

	rootCmd.PersistentFlags().StringVarP(&network, "network", "n", "simnet", "specify XUD network")
	rootCmd.PersistentFlags().StringVar(&simnetDir, "simnet-dir", "", "Simnet environment folder")
	rootCmd.PersistentFlags().StringVar(&testnetDir, "testnet-dir", "", "Testnet environment folder")
	rootCmd.PersistentFlags().StringVar(&mainnetDir, "mainnet-dir", "", "Mainnet environment folder")

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
		networkDir = filepath.Join(homeDir, "simnet")
		if viper.GetString("simnet-dir") != "" {
			networkDir = viper.GetString("simnet-dir")
		}
		if simnetDir != "" {
			networkDir = simnetDir
		}
	case "testnet":
		networkDir = filepath.Join(homeDir, "testnet")
		if viper.GetString("testnet-dir") != "" {
			networkDir = viper.GetString("testnet-dir")
		}
		if testnetDir != "" {
			networkDir = testnetDir
		}
	case "mainnet":
		networkDir = filepath.Join(homeDir, "mainnet")
		if viper.GetString("mainnet-dir") != "" {
			networkDir = viper.GetString("mainnet-dir")
		}
		if mainnetDir != "" {
			networkDir = mainnetDir
		}
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
}
