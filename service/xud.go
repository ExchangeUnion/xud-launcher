package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type XudConfig struct {
	PreserveConfig string `usage:"Preserve xud xud.conf file during updates"`
}

type Xud struct {
	Base

	Config XudConfig

	PreserveConfig bool
}

func newXud(name string) Xud {
	return Xud{
		Base: newBase(name),
	}
}

func (t *Xud) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network, &BaseConfig{
		Disabled:    false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       images[network][t.Name],
	}); err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.Config, &XudConfig{
		PreserveConfig: "",
	}, cmd); err != nil {
		return err
	}

	return nil
}

func (t *Xud) GetConfig() interface{} {
	return t.Config
}

func (t *Xud) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.xud", config.Network)
	if err != nil {
		return err
	}

	// xud apply
	t.Environment["NETWORK"] = network
	t.Environment["NODE_ENV"] = "production"

	if t.Config.PreserveConfig == "true" {
		t.Environment["PRESERVE_CONFIG"] = "true"
	} else {
		t.Environment["PRESERVE_CONFIG"] = "false"
	}

	t.Volumes = append(t.Volumes,
		"./data/lndbtc:/root/.lndbtc",
		"./data/lndltc:/root/.lndltc",
		//"/:/mnt/hostfs",
		"./backup:/root/backup",
	)

	switch network {
	case "simnet":
		t.Ports = append(t.Ports, "28885")
	case "testnet":
		t.Ports = append(t.Ports, "18885")
	case "mainnet":
		t.Ports = append(t.Ports, "8885")
	}

	return nil
}

func (t *Xud) ToJson() map[string]interface{} {
	result := t.Base.ToJson()

	rpc := make(map[string]interface{})
	result["rpc"] = rpc
	rpc["type"] = "gRPC"
	rpc["host"] = "xud"
	switch t.Network {
	case "simnet":
		rpc["port"] = 28886
	case "testnet":
		rpc["port"] = 18886
	case "mainnet":
		rpc["port"] = 8886
	}
	rpc["tlsCert"] = fmt.Sprintf("%s/xud/tls.cert", PROXY_DATA_DIR)

	return result
}
