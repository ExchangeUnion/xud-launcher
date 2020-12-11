package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

// TODO support variable replacement in usage tag
type LndConfig struct {
	Mode           string `usage:"Lnd service mode"`
	PreserveConfig string `usage:"Preserve lnd.conf file during updates"`
}

type Lnd struct {
	Base

	Config LndConfig
	Chain  string
}

func newLnd(name string, chain string) Lnd {
	return Lnd{
		Base:  newBase(name),
		Chain: chain,
	}
}

func (t *Lnd) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network, &BaseConfig{
		Disabled:    false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       images[network][t.Name],
	}); err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.Config, &LndConfig{
		Mode:           "native",
		PreserveConfig: "false",
	}, cmd); err != nil {
		return err
	}

	return nil
}

func (t *Lnd) GetConfig() interface{} {
	return t.Config
}

func (t *Lnd) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}
	t.Network = network

	// base apply
	err := t.Base.Apply("/root/.lnd", config.Network)
	if err != nil {
		return err
	}

	// lnd apply
	t.Environment["NETWORK"] = network
	t.Environment["CHAIN"] = t.Chain

	if t.Config.PreserveConfig == "true" {
		t.Environment["PRESERVE_CONFIG"] = "true"
	} else {
		t.Environment["PRESERVE_CONFIG"] = "false"
	}

	if config.ExternalIp != "" {
		t.Environment["EXTERNAL_IP"] = config.ExternalIp
	}

	if network == "testnet" || network == "mainnet" {
		var mode string
		var rpchost string
		var rpcport uint16
		var rpcuser string
		var rpcpass string
		var zmqpubrawblock string
		var zmqpubrawtx string

		if t.Chain == "bitcoin" {
			backend := services["bitcoind"].GetConfig().(BitcoindConfig)
			mode = backend.Mode
			rpchost = backend.Rpchost
			rpcport = backend.Rpcport
			rpcuser = backend.Rpcuser
			rpcpass = backend.Rpcpass
			zmqpubrawblock = backend.Zmqpubrawblock
			zmqpubrawtx = backend.Zmqpubrawtx
		} else {
			backend := services["litecoind"].GetConfig().(LitecoindConfig)
			mode = backend.Mode
			rpchost = backend.Rpchost
			rpcport = backend.Rpcport
			rpcuser = backend.Rpcuser
			rpcpass = backend.Rpcpass
			zmqpubrawblock = backend.Zmqpubrawblock
			zmqpubrawtx = backend.Zmqpubrawtx
		}

		if mode == "neutrino" || mode == "light" {
			t.Environment["NEUTRINO"] = "True"
		} else if mode == "external" {
			t.Environment["RPCHOST"] = rpchost
			t.Environment["RPCPORT"] = fmt.Sprint(rpcport)
			t.Environment["RPCUSER"] = rpcuser
			t.Environment["RPCPASS"] = rpcpass
			t.Environment["ZMQPUBRAWBLOCK"] = zmqpubrawblock
			t.Environment["ZMQPUBRAWTX"] = zmqpubrawtx
		}
	}

	if network == "simnet" {
		if t.Chain == "bitcoin" {
			t.Command = append(t.Command,
				"--debuglevel=debug",
				"--nobootstrap",
				"--minbackoff=30s",
				"--maxbackoff=24h",
				"--bitcoin.active",
				"--bitcoin.simnet",
				"--bitcoin.node=neutrino",
				"--bitcoin.defaultchanconfs=6",
				"--routing.assumechanvalid",
				"--neutrino.connect=btcd.simnet.exchangeunion.com:38555",
				"--chan-enable-timeout=0m10s",
				"--max-cltv-expiry=5000",
			)
		} else { // litecoin
			t.Command = append(t.Command,
				"--debuglevel=debug",
				"--nobootstrap",
				"--minbackoff=30s",
				"--maxbackoff=24h",
				"--litecoin.active",
				"--litecoin.simnet",
				"--litecoin.node=neutrino",
				"--litecoin.defaultchanconfs=6",
				"--routing.assumechanvalid",
				"--neutrino.connect=btcd.simnet.exchangeunion.com:39555",
				"--chan-enable-timeout=0m10s",
				"--max-cltv-expiry=20000",
			)
		}
	}

	t.Hostname = t.Name

	return nil
}

func (t *Lnd) ToJson() map[string]interface{} {
	result := t.Base.ToJson()

	result["mode"] = t.Config.Mode

	rpc := make(map[string]interface{})
	result["rpc"] = rpc
	rpc["type"] = "gRPC"

	var name string
	switch t.Chain {
	case "bitcoin":
		name = "lndbtc"
	case "litecoin":
		name = "lndltc"
	}
	rpc["host"] = name
	rpc["port"] = 10009
	rpc["tlsCert"] = fmt.Sprintf("%s/%s/tls.cert", PROXY_DATA_DIR, name)
	rpc["macaroon"] = fmt.Sprintf("%s/%s/data/chain/%s/%s/readonly.macaroon", PROXY_DATA_DIR, name, t.Chain, t.Network)

	return result
}
