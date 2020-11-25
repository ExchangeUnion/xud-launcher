package service

import (
	"fmt"
	"github.com/spf13/cobra"
)

type BaseConfig struct {
	Disable     bool
	ExposePorts []string
}

func configureCommonFlags(service string, config *BaseConfig, cmd *cobra.Command) error {
	cmd.PersistentFlags().BoolVar(
		&config.Disable,
		fmt.Sprintf("%s.disabled", service),
		false,
		fmt.Sprintf("Enable/Disable %s service", service),
	)
	cmd.PersistentFlags().StringSliceVar(
		&config.ExposePorts,
		fmt.Sprintf("%s.expose-ports", service),
		[]string{},
		fmt.Sprintf("Expose %s service ports to your host machine", service),
	)

	return nil
}

type Base struct {
	Image string
	Environment map[string]string
	Command []string
	Ports []string
	Volumes []string
}

type Service interface {
	Apply(network string) error
}
