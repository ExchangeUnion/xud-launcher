package service

import (
	"errors"
	"github.com/spf13/cobra"
)

type XudConfig struct {
	PreserveConfig bool
}

type Xud struct {
	Base

	config XudConfig
}

func newXud(name string) Xud {
	return Xud{
		Base: newBase(name),
	}
}

func (t Xud) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         "./data/xud",
	}, cmd)
	if err != nil {
		return err
	}

	cmd.PersistentFlags().BoolVar(&t.config.PreserveConfig, "xud.preserve-config", false, "Preserve xud xud.conf file during updates")

	return nil
}

func (t Xud) GetConfig() interface{} {
	return t.config
}

func (t Xud) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.xud")
	if err != nil {
		return err
	}

	// xud apply
	t.Environment["NETWORK"] = network
	t.Environment["NODE_ENV"] = "production"

	if t.config.PreserveConfig {
		t.Environment["PRESERVE_CONFIG"] = "true"
	} else {
		t.Environment["PRESERVE_CONFIG"] = "false"
	}

	return nil
}
