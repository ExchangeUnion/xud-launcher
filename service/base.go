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

func configureCommonFlags(service string, config *BaseConfig, defaultValues *BaseConfig, cmd *cobra.Command) error {
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

type Base struct {
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

type Service interface {
	GetConfig() interface{}
	Apply(network string, services map[string]Service) error
}
