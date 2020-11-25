package service

import (
	"fmt"
	"github.com/spf13/cobra"
)

type BaseConfig struct {
	Disable     bool
	ExposePorts []string
	Dir         string
}

func configureBaseFlags(service string, config *BaseConfig, defaultValues *BaseConfig, cmd *cobra.Command) error {
	cmd.PersistentFlags().BoolVar(
		&config.Disable,
		fmt.Sprintf("%s.disabled", service),
		defaultValues.Disable,
		fmt.Sprintf("Enable/Disable %s service", service),
	)
	cmd.PersistentFlags().StringSliceVar(
		&config.ExposePorts,
		fmt.Sprintf("%s.expose-ports", service),
		defaultValues.ExposePorts,
		fmt.Sprintf("Expose %s service ports to your host machine", service),
	)
	cmd.PersistentFlags().StringVar(
		&config.Dir,
		fmt.Sprintf("%s.dir", service),
		defaultValues.Dir,
		fmt.Sprintf("Specify the main data directory of %s service", service),
	)

	return nil
}

type SharedConfig struct {
	Network string
	SimnetDir string
	TestnetDir string
	MainnetDir string
	ExternalIp string
	Dev bool
	UseLocalImages string
}

type Base struct {
	Name string
	Image       string
	Environment map[string]string
	Command     []string
	Ports       []string
	Volumes     []string
}

func (t Base) Apply(config *BaseConfig, dir string, network string, services map[string] Service) error {
	for _, port := range config.ExposePorts {
		t.Ports = append(t.Ports, port)
	}

	t.Volumes = append(t.Volumes, fmt.Sprintf("%s:%s", config.Dir, dir))

	return nil
}

func (t Base) GetName() string {
	return t.Name
}

type Service interface {
	ConfigureFlags(cmd *cobra.Command) error
	GetConfig() interface{}
	GetName() string
	Apply(config *SharedConfig, services map[string]Service) error
}

func NewService(name string) Service {
	if name == "bitcoind" {
		return newBitcoind("bitcoind")
	} else if name == "litecoind" {
		return newLitecoind("litecoind")
	} else if name == "geth" {
		return newGeth("geth")
	} else if name == "lndbtc" {
		return newLnd("lndbtc", "bitcoin")
	} else if name == "lndltc" {
		return newLnd("lndltc","litecoin")
	} else if name == "connext" {
		return newConnext("connext")
	} else if name == "xud" {
		return newXud("xud")
	} else if name == "arby" {
		return newArby("arby")
	} else if name == "boltz" {
		return newBoltz("boltz")
	} else if name == "webui" {
		return newWebui("webui")
	} else if name == "proxy" {
		return newProxy("proxy")
	}

	return nil
}
