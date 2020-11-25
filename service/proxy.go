package service

import (
	"errors"
	"github.com/spf13/cobra"
)

type ProxyConfig struct {
	BaseConfig

	// add more proxy specified attributes here
}

type Proxy struct {
	Base

	config ProxyConfig
}

func newProxy(name string) Proxy {
	return Proxy{
		Base: Base{
			Name: name,
		},
	}
}

func (t Proxy) ConfigureFlags(cmd *cobra.Command) error {
	err := configureBaseFlags(t.Name, &t.config.BaseConfig, &BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         "./data/proxy",
	}, cmd)
	if err != nil {
		return err
	}

	// configure proxy specified flags here

	return nil
}

func (t Proxy) GetConfig() interface{} {
	return t.config
}

func (t Proxy) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply(&t.config.BaseConfig, "/root/.proxy", network, services)
	if err != nil {
		return err
	}

	// proxy apply

	return nil
}
