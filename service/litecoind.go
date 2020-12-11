package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type LitecoindConfig struct {
	Mode           string `usage:"Litecoind service mode"`
	Rpchost        string `usage:"External litecoind RPC hostname"`
	Rpcport        uint16 `usage:"External litecoind RPC port"`
	Rpcuser        string `usage:"External litecoind RPC username"`
	Rpcpass        string `usage:"External litecoind RPC password"`
	Zmqpubrawblock string `usage:"External litecoind ZeroMQ raw blocks publication address"`
	Zmqpubrawtx    string `usage:"External litecoind ZeroMQ raw transactions publication address"`
}

type Litecoind struct {
	Base

	Config LitecoindConfig
}

func newLitecoind(name string) Litecoind {
	return Litecoind{
		Base: newBase(name),
	}
}

func (t *Litecoind) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network, &BaseConfig{
		Disabled:    true,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       images[network][t.Name],
	}); err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.Config, &LitecoindConfig{
		Mode:           "light",
		Rpchost:        "",
		Rpcport:        0,
		Rpcuser:        "",
		Rpcpass:        "",
		Zmqpubrawblock: "",
		Zmqpubrawtx:    "",
	}, cmd); err != nil {
		return err
	}

	return nil
}

func (t *Litecoind) GetConfig() interface{} {
	return t.Config
}

func (t *Litecoind) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.litecoind", config.Network)
	if err != nil {
		return err
	}

	// litecoind apply
	t.Environment["NETWORK"] = network

	if t.Config.Mode != "native" || network == "simnet" {
		t.Disabled = true
	}

	return nil
}

func (t *Litecoind) ToJson() map[string]interface{} {
	result := t.Base.ToJson()
	result["mode"] = t.Config.Mode

	rpc := make(map[string]interface{})
	result["rpc"] = rpc
	rpc["type"] = "JSON-RPC"
	rpc["host"] = "litecoind"
	if t.Network == "testnet" {
		rpc["port"] = 19332
	} else {
		rpc["port"] = 9332
	}
	rpc["username"] = "xu"
	rpc["password"] = "xu"

	return result
}
