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

	Config BoltzConfig
}

func newBoltz(name string) Boltz {
	return Boltz{
		Base: newBase(name),
	}
}

func (t *Boltz) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network); err != nil {
		return err
	}

	// configure boltz specified flags here

	return nil
}

func (t *Boltz) GetConfig() interface{} {
	return t.Config
}

func (t *Boltz) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.boltz", config.Network)
	if err != nil {
		return err
	}

	// boltz apply

	if network == "simnet" {
		t.Disabled = true
	}

	return nil
}

func (t *Boltz) ToJson() map[string]interface{} {
	result := t.Base.ToJson()

	rpc := make(map[string]interface{})
	result["rpc"] = rpc
	rpc["type"] = "gRPC"
	rpc["host"] = "boltz"
	rpc["btcPort"] = 9002
	rpc["ltcPort"] = 9003

	return result
}
