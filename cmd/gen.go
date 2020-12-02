package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/reliveyy/xud-launcher/service"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	config   service.SharedConfig
	services map[string]service.Service
)

func init() {
	genCmd.PersistentFlags().StringVar(&config.ExternalIp, "external-ip", "", "Host machine external IP address")
	genCmd.PersistentFlags().BoolVar(&config.Dev, "dev", false, "Use local built utils image")
	genCmd.PersistentFlags().StringVar(&config.UseLocalImages, "use-local-images", "", "Use other local built images")

	// [Add capability to restrict flag values to a set of allowed values](https://github.com/spf13/pflag/issues/236)
	names := []string{
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

	logger.Info("Configuring subcommand flags")

	for _, name := range names {
		s := service.NewService(name)
		err := s.ConfigureFlags(genCmd)
		if err != nil {
			logger.Fatal(err)
		}
		services[name] = s
	}

	rootCmd.AddCommand(genCmd)
}

func Export(services []service.Service) string {
	var result = ""

	result += "version: \"2.4\"\n"

	result += "services:\n"

	for _, s := range services {

		result += fmt.Sprintf("  %s:\n", s.GetName())

		result += fmt.Sprintf("    image: %s\n", s.GetImage())

		if s.GetHostname() != "" {
			result += fmt.Sprintf("    hostname: %s\n", s.GetName())
		}

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

		if s.GetName() == "xud" {
			result += fmt.Sprintf("    entrypoint: [\"bash\", \"-c\", \"echo /root/backup > /root/.xud/.backup-dir-value && /entrypoint.sh\"]\n")
		}
	}

	return result
}

func ExportCompose(services []service.Service) {
	path := filepath.Join(networkDir, "docker-compose.yml")
	f, err := os.Create(path)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	var targets []service.Service

	// filter enabled services
	for _, s := range services {
		if s.IsDisabled() {
			continue
		}
		targets = append(targets, s)
	}

	yml := Export(targets)
	_, err = f.WriteString(yml)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Exported to %s\n", path)
}

type Service struct {
	Name     string      `json:"name"`
	Rpc      interface{} `json:"rpc"`
	Disabled bool        `json:"disabled"`
	Mode     string      `json:"mode"`
}

type Config struct {
	Timestamp string        `json:"timestamp"`
	Network   string        `json:"network"`
	Services  []interface{} `json:"services"`
}

func Export2(services []service.Service) string {
	var config Config
	config.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	config.Network = network
	for _, s := range services {
		if s.GetName() == "proxy" {
			continue
		}
		config.Services = append(config.Services, s.ToJson())
	}
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		logger.Fatal(err)
	}
	return string(data)
}

func ExportConfig(services []service.Service) {
	path := filepath.Join(dataDir, "config.json")
	f, err := os.Create(path)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	j := Export2(services)
	_, err = f.WriteString(j)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Exported to %s\n", path)
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate docker-compose.yml file from xud-docker configurations",
	Run: func(cmd *cobra.Command, args []string) {
		config.Network = network
		config.HomeDir = homeDir
		config.NetworkDir = networkDir

		var order []string
		if network == "simnet" {
			order = []string{
				"lndbtc",
				"lndltc",
				"connext",
				"xud",
				"arby",
				"webui",
				"proxy",
			}
		} else {
			order = []string{
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
		}

		var targets []service.Service

		// apply for all services
		for _, name := range order {
			s := services[name]
			err := s.Apply(&config, services)
			if err != nil {
				logger.Fatalf("%s: %s", name, err)
			}
			targets = append(targets, s)
		}

		ExportCompose(targets)
		ExportConfig(targets)
	},
}
