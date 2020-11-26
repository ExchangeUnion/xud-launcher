package service

import (
	"errors"
	"github.com/spf13/cobra"
)

type BoltzConfig struct {
	// add more boltz specified attributes here
}

type Boltz struct {
	Base

	config BoltzConfig
}

func newBoltz(name string) Boltz {
	return Boltz{
		Base: newBase(name),
	}
}

func (t Boltz) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         "./data/boltz",
	}, cmd)
	if err != nil {
		return err
	}

	// configure boltz specified flags here

	return nil
}

func (t Boltz) GetConfig() interface{} {
	return t.config
}

func (t Boltz) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.boltz")
	if err != nil {
		return err
	}

	// boltz apply

	return nil
}
