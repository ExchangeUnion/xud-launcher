package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type GethConfig struct {
	Mode                string `usage:"Geth service mode"`
	Rpchost             string `usage:"External geth RPC hostname"`
	Rpcport             uint16 `usage:"External geth RPC port"`
	InfuraProjectId     string `usage:"Infura geth provider project ID"`
	InfuraProjectSecret string `usage:"Infura geth provider project secret"`
	Cache               string `usage:"Geth cache size"`
	AncientChaindataDir string `usage:"Specify the container's volume mapping ancient chaindata directory. Can be located on a slower HDD."`
}

type Geth struct {
	Base

	Config GethConfig
}

func newGeth(name string) Geth {
	return Geth{
		Base: newBase(name),
	}
}

func (t *Geth) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network); err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.Config, &GethConfig{
		Mode:                "light",
		Rpchost:             "",
		Rpcport:             0,
		InfuraProjectId:     "",
		InfuraProjectSecret: "",
		Cache:               "",
		AncientChaindataDir: "",
	}, cmd); err != nil {
		return err
	}

	return nil
}

func (t *Geth) GetConfig() interface{} {
	return t.Config
}

func (t *Geth) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.ethereum", config.Network)
	if err != nil {
		return err
	}

	// geth apply
	t.Environment["NETWORK"] = network

	if t.Config.AncientChaindataDir != "" {
		volume := fmt.Sprintf("%s:/root/.ethereum-ancient-chaindata", t.Config.AncientChaindataDir)
		t.Volumes = append(t.Volumes, volume)
		t.Environment["CUSTOM_ANCIENT_CHAINDATA"] = "true"
	}

	if t.Config.Cache != "" {
		t.Command = append(t.Command, fmt.Sprintf("--cache %s", t.Config.Cache))
	}

	// TODO select ethProvider in light mode

	if t.Config.Mode != "native" || network == "simnet" {
		t.Disabled = true
	}

	return nil
}

func (t *Geth) ToJson() map[string]interface{} {
	result := t.Base.ToJson()
	result["mode"] = t.Config.Mode

	rpc := make(map[string]interface{})
	result["rpc"] = rpc
	rpc["type"] = "JSON-RPC"
	rpc["host"] = "geth"
	rpc["port"] = 8545

	return result
}
