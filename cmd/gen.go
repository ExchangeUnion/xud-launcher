package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/reliveyy/xud-launcher/service"
	"github.com/spf13/cobra"
	"log"
	"path"
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
		generalConf := path.Join(homeDir, "xud-docker.conf")
		networkDir := path.Join(homeDir, network)
		networkConf := path.Join(networkDir, fmt.Sprintf("%s.conf", network))

		fmt.Println(generalConf)
		fmt.Println(networkConf)

		config.Network = network

		println("version: 2.4")
		println("services:")

		for _, name := range names {
			service := services[name]

			if network == "simnet" {
				if name == "bitcoind" || name == "litecoind" || name == "geth" || name == "boltz" {
					continue
				}
			}

			if service.Disabled() {
				continue
			}

			err = service.Apply(&config, services)
			if err != nil {
				log.Fatalf("%s: %s", name, err)
			}

			fmt.Printf("  %s:\n", name)

			fmt.Printf("    image: %s\n", service.GetImage())

			if len(service.GetCommand()) > 0 {
				fmt.Printf("    command: >\n")
				for _, arg := range service.GetCommand() {
					fmt.Printf("      %s\n", arg)
				}
			}

			if len(service.GetEnvironment()) > 0 {
				fmt.Printf("    environment:\n")
				for key, value := range service.GetEnvironment() {
					fmt.Printf("    - %s=%s\n", key, value)
				}
			}

			if len(service.GetVolumes()) > 0 {
				fmt.Printf("    volumes:\n")
				for _, volume := range service.GetVolumes() {
					fmt.Printf("    - %s\n", volume)
				}
			}

			if len(service.GetPorts()) > 0 {
				fmt.Printf("    ports:\n")
				for _, port := range service.GetPorts() {
					fmt.Printf("    - %s\n", port)
				}
			}

		}
	},
}
