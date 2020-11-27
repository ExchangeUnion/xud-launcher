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

	config LitecoindConfig
}

func newLitecoind(name string) Litecoind {
	return Litecoind{
		Base: newBase(name),
	}
}

func (t *Litecoind) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       "exchangeunion/litecoind",
	}, cmd)
	if err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.config, cmd); err != nil {
		return err
	}

	//cmd.PersistentFlags().StringVar(&t.config.Mode, "litecoind.mode", "light", "Litecoind service mode")
	//cmd.PersistentFlags().StringVar(&t.config.Rpchost, "litecoind.rpchost", "", "External litecoind RPC hostname")
	//cmd.PersistentFlags().Uint16Var(&t.config.Rpcport, "litecoind.rpcport", 0, "External litecoind RPC port")
	//cmd.PersistentFlags().StringVar(&t.config.Rpcuser, "litecoind.rpcuser", "", "External litecoind RPC username")
	//cmd.PersistentFlags().StringVar(&t.config.Rpcpass, "litecoind.rpcpass", "", "External litecoind RPC password")
	//cmd.PersistentFlags().StringVar(&t.config.Zmqpubrawblock, "litecoind.zmqpubrawblock", "", "External litecoind ZeroMQ raw blocks publication address")
	//cmd.PersistentFlags().StringVar(&t.config.Zmqpubrawtx, "litecoind.zmqpubrawtx", "", "External litecoind ZeroMQ raw transactions publication address")

	return nil
}

func (t *Litecoind) GetConfig() interface{} {
	return t.config
}

func (t *Litecoind) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.litecoind", network)
	if err != nil {
		return err
	}

	// litecoind apply
	t.Environment["NETWORK"] = network

	return nil
}
