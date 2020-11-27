package service

import (
	"errors"
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

	config LitecoindConfig
}

func newLitecoind(name string) Litecoind {
	return Litecoind{
		Base: newBase(name),
	}
}

func (t *Litecoind) ConfigureFlags(cmd *cobra.Command) error {
	if err := t.Base.ConfigureFlags(cmd, true); err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.config, &LitecoindConfig{
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
	return t.config
}

func (t *Litecoind) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.config)

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

	return nil
}
