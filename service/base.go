package service

import (
	"fmt"
	"github.com/spf13/cobra"
)

type BaseConfig struct {
	Disable     bool
	ExposePorts []string
	Dir         string
	Image       string
}

type SharedConfig struct {
	Network        string
	SimnetDir      string
	TestnetDir     string
	MainnetDir     string
	ExternalIp     string
	Dev            bool
	UseLocalImages string
}

type Base struct {
	Name        string
	Image       string
	Environment map[string]string
	Command     []string
	Ports       []string
	Volumes     []string

	config BaseConfig
}

func newBase(name string) Base {
	return Base{
		Name:        name,
		Image:       "",
		Environment: make(map[string]string),
		Command:     []string{},
		Ports:       []string{},
		Volumes:     []string{},
	}
}

func (t Base) ConfigureFlags(defaultValues *BaseConfig, cmd *cobra.Command) error {
	cmd.PersistentFlags().BoolVar(
		&t.config.Disable,
		fmt.Sprintf("%s.disabled", t.Name),
		defaultValues.Disable,
		fmt.Sprintf("Enable/Disable %s service", t.Name),
	)
	t.config.Disable = defaultValues.Disable

	cmd.PersistentFlags().StringSliceVar(
		&t.config.ExposePorts,
		fmt.Sprintf("%s.expose-ports", t.Name),
		defaultValues.ExposePorts,
		fmt.Sprintf("Expose %s service ports to your host machine", t.Name),
	)
	t.config.ExposePorts = defaultValues.ExposePorts

	cmd.PersistentFlags().StringVar(
		&t.config.Dir,
		fmt.Sprintf("%s.dir", t.Name),
		defaultValues.Dir,
		fmt.Sprintf("Specify the main data directory of %s service", t.Name),
	)
	t.config.Dir = defaultValues.Dir

	cmd.PersistentFlags().StringVar(
		&t.config.Image,
		fmt.Sprintf("%s.image", t.Name),
		defaultValues.Image,
		fmt.Sprintf("Specify the image of %s service", t.Name),
	)
	t.config.Image = defaultValues.Image

	return nil
}

func (t Base) Apply(dir string) error {
	for _, port := range t.config.ExposePorts {
		t.Ports = append(t.Ports, port)
	}

	t.Volumes = append(t.Volumes, fmt.Sprintf("%s:%s", t.config.Dir, dir))
	t.Image = t.config.Image

	return nil
}

func (t Base) GetName() string {
	return t.Name
}

func (t Base) GetImage() string {
	return t.Image
}

func (t Base) GetCommand() []string {
	return t.Command
}

func (t Base) GetEnvironment() map[string]string {
	return t.Environment
}

func (t Base) GetVolumes() []string {
	return t.Volumes
}

func (t Base) GetPorts() []string {
	return t.Ports
}

func (t Base) Disabled() bool {
	return t.config.Disable
}

type Service interface {
	ConfigureFlags(cmd *cobra.Command) error
	GetConfig() interface{}
	GetName() string
	Apply(config *SharedConfig, services map[string]Service) error

	GetImage() string
	GetCommand() []string
	GetEnvironment() map[string]string
	GetVolumes() []string
	GetPorts() []string
	Disabled() bool
}

func NewService(name string) Service {
	switch name {
	case "bitcoind":
		return newBitcoind("bitcoind")
	case "litecoind":
		return newLitecoind("litecoind")
	case "geth":
		return newGeth("geth")
	case "lndbtc":
		return newLnd("lndbtc", "bitcoin")
	case "lndltc":
		return newLnd("lndltc", "litecoin")
	case "connext":
		return newConnext("connext")
	case "xud":
		return newXud("xud")
	case "arby":
		return newArby("arby")
	case "boltz":
		return newBoltz("boltz")
	case "webui":
		return newWebui("webui")
	case "proxy":
		return newProxy("proxy")
	}

	return nil
}
