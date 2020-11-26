package service

import (
	"errors"
	"fmt"
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

func (t *Proxy) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       "exchangeunion/proxy",
	}, cmd)
	if err != nil {
		return err
	}

	// configure proxy specified flags here

	return nil
}

func (t *Proxy) GetConfig() interface{} {
	return t.config
}

func (t *Proxy) Apply(config *SharedConfig, services map[string]Service) error {
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
	t.Volumes = append(t.Volumes,
		"/var/run/docker.sock:/var/run/docker.sock",
		"./logs/config.sh:/root/config.sh",
		"./data/xud:/root/.xud",
		"./data/lndbtc:/root/.lndbtc",
		"./data/lndltc:/root/.lndltc",
	)

	switch network {
	case "simnet":
		t.Ports = append(t.Ports, "127.0.0.1:28889:8080")
	case "testnet":
		t.Ports = append(t.Ports, "127.0.0.1:18889:8080")
	case "mainnet":
		t.Ports = append(t.Ports, "127.0.0.1:8889:8080")
	}

	return nil
}
