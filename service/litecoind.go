package service

import (
	"errors"
	"github.com/spf13/cobra"
)

type LitecoindConfig struct {
	Mode           string
	Rpchost        string
	Rpcport        uint16
	Rpcuser        string
	Rpcpass        string
	Zmqpubrawblock string
	Zmqpubrawtx    string
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

func (t Litecoind) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         "./data/litecoind",
	}, cmd)
	if err != nil {
		return err
	}

	cmd.PersistentFlags().StringVar(&t.config.Mode, "litecoind.mode", "light", "Litecoind service mode")
	cmd.PersistentFlags().StringVar(&t.config.Rpchost, "litecoind.rpchost", "", "External litecoind RPC hostname")
	cmd.PersistentFlags().Uint16Var(&t.config.Rpcport, "litecoind.rpcport", 0, "External litecoind RPC port")
	cmd.PersistentFlags().StringVar(&t.config.Rpcuser, "litecoind.rpcuser", "", "External litecoind RPC username")
	cmd.PersistentFlags().StringVar(&t.config.Rpcpass, "litecoind.rpcpass", "", "External litecoind RPC password")
	cmd.PersistentFlags().StringVar(&t.config.Zmqpubrawblock, "litecoind.zmqpubrawblock", "", "External litecoind ZeroMQ raw blocks publication address")
	cmd.PersistentFlags().StringVar(&t.config.Zmqpubrawtx, "litecoind.zmqpubrawtx", "", "External litecoind ZeroMQ raw transactions publication address")

	return nil
}

func (t Litecoind) GetConfig() interface{} {
	return t.config
}

func (t Litecoind) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.litecoind")
	if err != nil {
		return err
	}

	// litecoind apply
	t.Environment["NETWORK"] = network

	return nil
}
