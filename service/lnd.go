package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type LndConfig struct {
	Mode           string
	PreserveConfig bool
}

type Lnd struct {
	Base

	config LndConfig
	Chain  string
}

func newLnd(name string, chain string) Lnd {
	return Lnd{
		Base:  newBase(name),
		Chain: chain,
	}
}

func (t *Lnd) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       fmt.Sprintf("exchangeunion/%s", t.Name),
	}, cmd)
	if err != nil {
		return err
	}

	cmd.PersistentFlags().StringVar(
		&t.config.Mode,
		fmt.Sprintf("%s.mode", t.Name),
		"native",
		fmt.Sprintf("%s service mode", strings.Title(t.Name)),
	)
	cmd.PersistentFlags().BoolVar(
		&t.config.PreserveConfig,
		fmt.Sprintf("%s.preserve-config", t.Name),
		false,
		fmt.Sprintf("Preserve %s lnd.conf file during updates", t.Name),
	)

	return nil
}

func (t *Lnd) GetConfig() interface{} {
	return t.config
}

func (t *Lnd) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.lnd", network)
	if err != nil {
		return err
	}

	// lnd apply
	t.Environment["NETWORK"] = network
	t.Environment["CHAIN"] = t.Chain

	if t.config.PreserveConfig {
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

	return nil
}
