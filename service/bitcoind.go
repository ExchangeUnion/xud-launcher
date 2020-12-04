package service

import (
	"errors"
	"github.com/spf13/cobra"
)

type BitcoindConfig struct {
	Mode           string `usage:"Bitcoind service mode"`
	Rpchost        string `usage:"External bitcoind RPC hostname"`
	Rpcport        uint16 `usage:"External bitcoind RPC port"`
	Rpcuser        string `usage:"External bitcoind RPC username"`
	Rpcpass        string `usage:"External bitcoind RPC password"`
	Zmqpubrawblock string `usage:"External bitcoind ZeroMQ raw blocks publication address"`
	Zmqpubrawtx    string `usage:"External bitcoind ZeroMQ raw transactions publication address"`
}

type Bitcoind struct {
	Base

	Config BitcoindConfig
}

func newBitcoind(name string) Bitcoind {
	return Bitcoind{
		Base: newBase(name),
	}
}

func (t *Bitcoind) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network); err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.Config, &BitcoindConfig{
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

func (t *Bitcoind) GetConfig() interface{} {
	return t.Config
}

func (t *Bitcoind) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply

	err := t.Base.Apply("/root/.bitcoind", config.Network)
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

	if t.Config.Mode != "native" || network == "simnet" {
		t.Disabled = true
	}

	return nil
}

func (t *Bitcoind) ToJson() map[string]interface{} {
	result := t.Base.ToJson()
	result["mode"] = t.Config.Mode

	rpc := make(map[string]interface{})
	result["rpc"] = rpc
	rpc["type"] = "JSON-RPC"
	rpc["host"] = "bitcoind"
	if t.Network == "testnet" {
		rpc["port"] = 18332
	} else {
		rpc["port"] = 8332
	}
	rpc["username"] = "xu"
	rpc["password"] = "xu"

	return result
}
