package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type XudConfig struct {
	PreserveConfig bool `usage:"Preserve xud xud.conf file during updates"`
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

func (t *Xud) ConfigureFlags(cmd *cobra.Command) error {
	if err := t.Base.ConfigureFlags(cmd, &BaseConfig{
		Disabled:    false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       "",
	}); err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.config, &XudConfig{
		PreserveConfig: false,
	}, cmd); err != nil {
		return err
	}

	return nil
}

func (t *Xud) GetConfig() interface{} {
	return t.config
}

func (t *Xud) Apply(config *SharedConfig, services map[string]Service) error {
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

	if t.config.PreserveConfig {
		t.Environment["PRESERVE_CONFIG"] = "true"
	} else {
		t.Environment["PRESERVE_CONFIG"] = "false"
	}

	t.Volumes = append(t.Volumes,
		"./data/lndbtc:/root/.lndbtc",
		"./data/lndltc:/root/.lndltc",
		"/:/mnt/hostfs",
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
