package cmd

import (
	"errors"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	logger        = initLogger()
	validNetworks = []string{"mainnet", "testnet", "simnet"}
	network       = initNetwork()

	homeDir    string
	networkDir string
	dataDir    string
	logsDir    string

	rootCmd = &cobra.Command{
		Use:   "xud-launcher",
		Short: fmt.Sprintf("XUD environment launcher (%s)", network),
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

func initNetwork() string {
	n := os.Getenv("NETWORK")
	n = strings.TrimSpace(n)
	n = strings.ToLower(n)
	if n == "" {
		logger.Debug("Use network: mainnet")
		return "mainnet" // default network
	}
	var valid = false
	for _, vn := range validNetworks {
		if n == vn {
			valid = true
			break
		}
	}
	if !valid {
		logger.Fatalf("Invalid network: %s", n)
	}
	logger.Debugf("Use network: %s", n)
	return n
}

func initLogger() *logrus.Entry {
	lv := os.Getenv("LOG_LEVEL")
	lv = strings.TrimSpace(lv)
	lv = strings.ToLower(lv)
	if lv == "" {
		logrus.SetLevel(logrus.WarnLevel)
	} else {
		switch lv {
		case "trace":
			logrus.SetLevel(logrus.TraceLevel)
		case "debug":
			logrus.SetLevel(logrus.DebugLevel)
		case "info":
			logrus.SetLevel(logrus.InfoLevel)
		case "warn":
			logrus.SetLevel(logrus.WarnLevel)
		case "error":
			logrus.SetLevel(logrus.ErrorLevel)
		case "fatal":
			logrus.SetLevel(logrus.FatalLevel)
		case "panic":
			logrus.SetLevel(logrus.PanicLevel)
		}
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true,
	})
	if runtime.GOOS == "windows" {
		logrus.SetOutput(colorable.NewColorableStdout())
	}
	logger := logrus.NewEntry(logrus.StandardLogger())
	return logger
}

func init() {
	homeDir = getHomeDir()
	ensureDir(homeDir)

	logger.Debugf("Configuring global flags")

	key := fmt.Sprintf("%s-dir", network)

	rootCmd.PersistentFlags().StringVar(
		&networkDir,
		key,
		filepath.Join(homeDir, network),
		fmt.Sprintf("%s environment folder", strings.Title(network)),
	)

	if err := viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(key)); err != nil {
		logger.Fatal("Failed to bind Viper key %s: %s", key, err)
	}

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	generalConf := filepath.Join(homeDir, "xud-docker.conf")
	logger.Debugf("Loading general config file: %s", generalConf)

	viper.SetConfigFile(generalConf)
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Debugf("Failed to load general config: %s", err)
	}

	ensureDir(networkDir)

	dataDir = filepath.Join(networkDir, "data")
	ensureDir(dataDir)
	logsDir = filepath.Join(networkDir, "logs")
	ensureDir(logsDir)

	logger.Debugf("homeDir=%s", homeDir)
	logger.Debugf("networkDir=%s", networkDir)
	logger.Debugf("dataDir=%s", dataDir)
	logger.Debugf("logsDir=%s", logsDir)

	networkConf := filepath.Join(networkDir, fmt.Sprintf("%s.conf", network))
	logger.Debugf("Loading network config file: %s", networkConf)

	viper.SetConfigFile(networkConf)
	viper.SetConfigType("toml")

	err = viper.MergeInConfig()
	if err != nil {
		logger.Debugf("Failed to load network config: %s", err)
	}
}
