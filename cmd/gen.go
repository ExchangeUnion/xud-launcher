package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/reliveyy/xud-launcher/service"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"strings"
)

var (
	config   service.SharedConfig
	services map[string]service.Service
	names    []string
)

func init() {
	genCmd.PersistentFlags().StringVar(&config.SimnetDir, "simnet-dir", "", "Simnet environment folder")
	genCmd.PersistentFlags().StringVar(&config.TestnetDir, "testnet-dir", "", "Testnet environment folder")
	genCmd.PersistentFlags().StringVar(&config.MainnetDir, "mainnet-dir", "", "Mainnet environment folder")
	genCmd.PersistentFlags().StringVar(&config.ExternalIp, "external-ip", "", "Host machine external IP address")
	genCmd.PersistentFlags().BoolVar(&config.Dev, "dev", false, "Use local built utils image")
	genCmd.PersistentFlags().StringVar(&config.UseLocalImages, "use-local-images", "", "Use other local built images")

	// [Add capability to restrict flag values to a set of allowed values](https://github.com/spf13/pflag/issues/236)
	names = []string{
		"bitcoind",
		"litecoind",
		"geth",
		"lndbtc",
		"lndltc",
		"connext",
		"xud",
		"arby",
		"boltz",
		"webui",
		"proxy",
	}

	services = make(map[string]service.Service)

	for _, name := range names {
		s := service.NewService(name)
		err := s.ConfigureFlags(genCmd)
		if err != nil {
			log.Fatal(err)
		}
		services[name] = s
	}

	rootCmd.AddCommand(genCmd)
}

func bypass(network string, s service.Service) bool {
	name := s.GetName()

	if network == "simnet" {
		if name == "bitcoind" || name == "litecoind" || name == "geth" || name == "boltz" {
			return true
		}
	}

	if s.Disabled() {
		return true
	}

	return false
}

func Export(services []service.Service) string {
	var result = ""

	result += "version: \"2.4\"\n"

	result += "services:\n"

	for _, s := range services {

		result += fmt.Sprintf("  %s:\n", s.GetName())

		result += fmt.Sprintf("    image: %s\n", s.GetImage())

		if len(s.GetCommand()) > 0 {
			result += fmt.Sprintf("    command: >\n")
			for _, arg := range s.GetCommand() {
				result += fmt.Sprintf("      %s\n", arg)
			}
		}

		if len(s.GetEnvironment()) > 0 {
			result += fmt.Sprintf("    environment:\n")
			for key, value := range s.GetEnvironment() {
				if strings.Contains(value, "\n") {
					// multiline value
					result += fmt.Sprintf("      - >\n")
					result += fmt.Sprintf("        %s=\n", key)
					for _, line := range strings.Split(value, "\n") {
						result += fmt.Sprintf("        %s\n", line)
					}
				} else {
					result += fmt.Sprintf("      - %s=%s\n", key, value)
				}

			}
		}

		if len(s.GetVolumes()) > 0 {
			result += fmt.Sprintf("    volumes:\n")
			for _, volume := range s.GetVolumes() {
				result += fmt.Sprintf("      - %s\n", volume)
			}
		}

		if len(s.GetPorts()) > 0 {
			result += fmt.Sprintf("    ports:\n")
			for _, port := range s.GetPorts() {
				result += fmt.Sprintf("      - %s\n", port)
			}
		}
	}

	return result
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate docker-compose.yml file from xud-docker configurations",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		homeDir = path.Join(homeDir, ".xud-docker")
		//generalConf := path.Join(homeDir, "xud-docker.conf")
		networkDir := path.Join(homeDir, network)
		//networkConf := path.Join(networkDir, fmt.Sprintf("%s.conf", network))

		config.Network = network

		var validServices = []service.Service{}

		for name, s := range services {
			if bypass(network, s) {
				continue
			}
			err = s.Apply(&config, services)
			if err != nil {
				log.Fatalf("%s: %s", name, err)
			}
			validServices = append(validServices, s)
		}

		composeFile := path.Join(networkDir, "docker-compose.yml")
		f, err := os.Create(composeFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		yml := Export(validServices)
		f.WriteString(yml)

		fmt.Printf("Exported to %s\n", composeFile)
	},
}
