package service

import (
	"errors"
	"github.com/spf13/cobra"
)

type ConnextConfig struct {
	// add more connext specified attributes here
}

type Connext struct {
	Base

	config ConnextConfig
}

func newConnext(name string) Connext {
	return Connext{
		Base: newBase(name),
	}
}

func (t Connext) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         "./data/connext",
	}, cmd)
	if err != nil {
		return err
	}

	// configure connext specified flags here

	return nil
}

func (t Connext) GetConfig() interface{} {
	return t.config
}

func (t Connext) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/app/connext-store")
	if err != nil {
		return err
	}

	// connext apply
	t.Environment["NETWORK"] = network
	t.Environment["VECTOR_CONFIG"] = `\
{
	"adminToken": "ddrWR8TK8UMTyR",
	"chainAddresses": {
		"1337": {
		"channelFactoryAddress": "0x2eC39861B9Be41c20675a1b727983E3F3151C576",
		"channelMastercopyAddress": "0x7AcAcA3BC812Bcc0185Fa63faF7fE06504D7Fa70",
		"transferRegistryAddress": "0xB2b8A1d98bdD5e7A94B3798A13A94dEFFB1Fe709",
		"TestToken": ""
		}
	},
	"chainProviders": {
		"1337": "http://35.234.110.95:8545"
	},
	"domainName": "",
	"logLevel": "debug",
	"messagingUrl": "https://messaging.connext.network",
	"production": true,
	"mnemonic": "crazy angry east hood fiber awake leg knife entire excite output scheme"
}
`
	t.Environment["VECTOR_SQLITE_FILE"] = "/database/store.db"
	t.Environment["VECTOR_PROD"] = "true"

	return nil
}
