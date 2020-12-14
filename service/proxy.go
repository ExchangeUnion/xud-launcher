package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

type ProxyConfig struct {
	// add more proxy specified attributes here
}

type Proxy struct {
	Base

	Config ProxyConfig
}

func newProxy(name string) Proxy {
	return Proxy{
		Base: newBase(name),
	}
}

func (t *Proxy) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network, &BaseConfig{
		Disabled:    false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       images[network][t.Name],
	}); err != nil {
		return err
	}

	// configure proxy specified flags here

	return nil
}

func (t *Proxy) GetConfig() interface{} {
	return t.Config
}

func (t *Proxy) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.proxy", config.Network)
	if err != nil {
		return err
	}

	var dockerSock string

	if runtime.GOOS == "windows" {
		dockerSock = "//var/run/docker.sock"
	} else {
		dockerSock = "/var/run/docker.sock"
	}

	// proxy apply
	t.Volumes = append(t.Volumes,
		fmt.Sprintf("%s:/var/run/docker.sock", dockerSock),
		fmt.Sprintf("%s:/root/network:ro", config.NetworkDir),
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
