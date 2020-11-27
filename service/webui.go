package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type WebuiConfig struct {
	// add more webui specified attributes here
}

type Webui struct {
	Base

	config WebuiConfig
}

func newWebui(name string) Webui {
	return Webui{
		Base: newBase(name),
	}
}

func (t *Webui) ConfigureFlags(cmd *cobra.Command) error {
	if err := t.Base.ConfigureFlags(cmd, &BaseConfig{
		Disabled:    true,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       "",
	}); err != nil {
		return err
	}

	// configure webui specified flags here

	return nil
}

func (t *Webui) GetConfig() interface{} {
	return t.config
}

func (t *Webui) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.webui", config.Network)
	if err != nil {
		return err
	}

	// webui apply

	return nil
}
