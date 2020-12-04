package service

import "github.com/spf13/cobra"

type Service interface {
	ConfigureFlags(cmd *cobra.Command, network string) error
	GetConfig() interface{}
	GetName() string
	Apply(config *SharedConfig, services map[string]Service) error

	GetImage() string
	GetCommand() []string
	GetEnvironment() map[string]string
	GetVolumes() []string
	GetPorts() []string
	GetHostname() string
	IsDisabled() bool

	ToJson() map[string]interface{}
}
