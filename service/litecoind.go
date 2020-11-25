package service

import (
	"errors"
	"github.com/spf13/cobra"
)

type LitecoindConfig struct {
	BaseConfig

	Mode string
	Rpchost string
	Rpcport uint16
	Rpcuser string
	Rpcpass string
	Zmqpubrawblock string
	Zmqpubrawtx string
}

type Litecoind struct {
	Base

	config LitecoindConfig
}

func NewLitecoind() Litecoind {
	return Litecoind{
		config: LitecoindConfig{},
	}
}

func (t Litecoind) ConfigureFlags(cmd *cobra.Command) error {
	err := configureCommonFlags("arby", &t.config.BaseConfig, cmd)
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

func (t Litecoind) Apply(network string) error {

	if network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	t.Environment["NETWORK"] = network

	return nil
}
