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
	branch string
	disableUpdate bool
	simnetDir string
	testnetDir string
	mainnetDir string
	externalIp string
	dev bool
	useLocalImages string
	api bool
)

func init() {
	genCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "master", "Git branch name")
	genCmd.PersistentFlags().BoolVar(&disableUpdate, "disable-update", false, "Skip update checks and enter xud-ctl shell directly")
	genCmd.PersistentFlags().StringVar(&simnetDir, "simnet-dir", "", "Simnet environment folder")
	genCmd.PersistentFlags().StringVar(&testnetDir, "testnet-dir", "", "Testnet environment folder")
	genCmd.PersistentFlags().StringVar(&mainnetDir, "mainnet-dir", "", "Mainnet environment folder")
	genCmd.PersistentFlags().StringVar(&externalIp, "external-ip", "", "Host machine external IP address")
	genCmd.PersistentFlags().BoolVar(&dev, "dev", false, "Use local built utils image")
	genCmd.PersistentFlags().StringVar(&useLocalImages, "use-local-images", "", "Use other local built images")
	genCmd.PersistentFlags().BoolVar(&api, "api", false, "Expose xud-docker API (REST + WebSocket)")

	// [Add capability to restrict flag values to a set of allowed values](https://github.com/spf13/pflag/issues/236)
	services := []string{
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

	for _, name := range services {
		s := service.NewService(name)
		err := s.ConfigureFlags(genCmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	genCmd.PersistentFlags().Bool("boltz.disabled", true, "Enable/Disable boltz service")

	genCmd.PersistentFlags().Bool("webui.disabled", true, "Enable/Disable webui service")
	genCmd.PersistentFlags().String("webui.expose-ports", "", "Expose webui service ports to your host machine")

	genCmd.PersistentFlags().Bool("proxy.disabled", true, "Enable/Disable proxy service")
	genCmd.PersistentFlags().String("proxy.expose-ports", "", "Expose proxy service ports to your host machine")
	
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
	},
}
