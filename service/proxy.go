package service

import (
	"errors"
	"github.com/spf13/cobra"
)

type ProxyConfig struct {
	// add more proxy specified attributes here
}

type Proxy struct {
	Base

	config ProxyConfig
}

func newProxy(name string) Proxy {
	return Proxy{
		Base: newBase(name),
	}
}

func (t Proxy) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
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
	err := t.Base.Apply("/root/.proxy")
	if err != nil {
		return err
	}

	// proxy apply

	return nil
}
