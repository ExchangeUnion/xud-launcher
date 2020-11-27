package cmd

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path"
)

var (
	// Used for flags.
	homeDir    string
	simnetDir  string
	testnetDir string
	mainnetDir string
	networkDir string
	network    string
	logger     = logrus.New()

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
	return path.Join(homeDir, ".xud-docker")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	logger.SetLevel(logrus.DebugLevel)

	homeDir = getHomeDir()

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
	generalConf := path.Join(homeDir, "xud-docker.conf")
	viper.SetConfigFile(generalConf)
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	logger.Infof("Loading general config file: %s", generalConf)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(err)
	}

	switch network {
	case "simnet":
		networkDir = path.Join(homeDir, "simnet")
		if viper.GetString("simnet-dir") != "" {
			networkDir = viper.GetString("simnet-dir")
		}
		if simnetDir != "" {
			networkDir = simnetDir
		}
	case "testnet":
		networkDir = path.Join(homeDir, "testnet")
		if viper.GetString("testnet-dir") != "" {
			networkDir = viper.GetString("testnet-dir")
		}
		if testnetDir != "" {
			networkDir = testnetDir
		}
	case "mainnet":
		networkDir = path.Join(homeDir, "mainnet")
		if viper.GetString("mainnet-dir") != "" {
			networkDir = viper.GetString("mainnet-dir")
		}
		if mainnetDir != "" {
			networkDir = mainnetDir
		}
	default:
		panic(errors.New("invalid network: " + network))
	}

	logger.Debugf("homeDir=%s", homeDir)
	logger.Debugf("networkDir=%s", networkDir)
}
