package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type BitcoindConfig struct {
	Mode           string
	Rpchost        string
	Rpcport        uint16
	Rpcuser        string
	Rpcpass        string
	Zmqpubrawblock string
	Zmqpubrawtx    string
}

type Bitcoind struct {
	Base

	config BitcoindConfig
}

func newBitcoind(name string) Bitcoind {
	return Bitcoind{
		Base: newBase(name),
	}
}

func (t *Bitcoind) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       "exchangeunion/bitcoind",
	}, cmd)
	if err != nil {
		return err
	}

	cmd.PersistentFlags().StringVar(&t.config.Mode, "bitcoind.mode", "light", "Bitcoind service mode")
	cmd.PersistentFlags().StringVar(&t.config.Rpchost, "bitcoind.rpchost", "", "External bitcoind RPC hostname")
	cmd.PersistentFlags().Uint16Var(&t.config.Rpcport, "bitcoind.rpcport", 0, "External bitcoind RPC port")
	cmd.PersistentFlags().StringVar(&t.config.Rpcuser, "bitcoind.rpcuser", "", "External bitcoind RPC username")
	cmd.PersistentFlags().StringVar(&t.config.Rpcpass, "bitcoind.rpcpass", "", "External bitcoind RPC password")
	cmd.PersistentFlags().StringVar(&t.config.Zmqpubrawblock, "bitcoind.zmqpubrawblock", "", "External bitcoind ZeroMQ raw blocks publication address")
	cmd.PersistentFlags().StringVar(&t.config.Zmqpubrawtx, "bitcoind.zmqpubrawtx", "", "External bitcoind ZeroMQ raw transactions publication address")

	return nil
}

func (t *Bitcoind) GetConfig() interface{} {
	return t.config
}

func (t *Bitcoind) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.bitcoind")
	if err != nil {
		return err
	}

	// bitcoind apply
	t.Environment["NETWORK"] = network

	t.Command = append(t.Command,
		"-server",
		"-rpcuser=xu",
		"-rpcpassword=xu",
		"-disablewallet",
		"-txindex",
		"-zmqpubrawblock=tcp://0.0.0.0:28332",
		"-zmqpubrawtx=tcp://0.0.0.0:28333",
		"-logips",
		"-rpcallowip=::/0",
		"-rpcbind=0.0.0.0",
	)

	if network == "testnet" {
		t.Command = append(t.Command, "-rpcport=18332", "-testnet")
	} else { // mainnet
		t.Command = append(t.Command, "-rpcport=8332")
	}

	return nil
}
