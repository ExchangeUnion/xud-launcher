package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type ArbyConfig struct {
	LiveCex                          bool   `usage:"Live CEX (deprecated)"`
	TestMode                         bool   `usage:"Whether to issue real orders on the centralized exchange"`
	BaseAsset                        string `usage:"Base asset"`
	QuoteAsset                       string `usage:"Quote asset"`
	CexBaseAsset                     string `usage:"Centralized exchange base asset"`
	CexQuoteAsset                    string `usage:"Centralized exchange quote asset"`
	TestCentralizedBaseassetBalance  string `usage:"Test centralized base asset balance"`
	TestCentralizedQuoteassetBalance string `usage:"Test centralized quote asset balance"`
	Cex                              string `usage:"Centralized Exchange"`
	CexApiKey                        string `usage:"CEX API key"`
	CexApiSecret                     string `usage:"CEX API secret"`
	Margin                           string `usage:"Trade margin"`
}

type Arby struct {
	Base

	Config ArbyConfig
}

func newArby(name string) Arby {
	return Arby{
		Base: newBase(name),
	}
}

func (t *Arby) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network); err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.Config, &ArbyConfig{
		LiveCex:                          true,
		TestMode:                         true,
		BaseAsset:                        "",
		QuoteAsset:                       "",
		CexBaseAsset:                     "",
		CexQuoteAsset:                    "",
		TestCentralizedBaseassetBalance:  "",
		TestCentralizedQuoteassetBalance: "",
		Cex:                              "binance",
		CexApiKey:                        "123",
		CexApiSecret:                     "abc",
		Margin:                           "0.04",
	}, cmd); err != nil {
		return err
	}

	if err := cmd.PersistentFlags().MarkDeprecated("arby.live-cex", "Please use --arby.test-mode instead"); err != nil {
		return err
	}

	return nil
}

func (t *Arby) GetConfig() interface{} {
	return t.Config
}

func (t *Arby) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.arby", config.Network)
	if err != nil {
		return err
	}

	// arby apply
	var rpcPort string

	t.Environment["NETWORK"] = network

	if network == "simnet" {
		rpcPort = "28886"
	} else if network == "testnet" {
		rpcPort = "18886"
	} else if network == "mainnet" {
		rpcPort = "8886"
	} else {
		return errors.New("invalid network: " + network)
	}

	t.Environment["NODE_ENV"] = "production"
	t.Environment["LOG_LEVEL"] = "trace"
	t.Environment["DATA_DIR"] = "/root/.arby"
	t.Environment["OPENDEX_CERT_PATH"] = "/root/.xud/tls.cert"
	t.Environment["OPENDEX_RPC_HOST"] = "xud"
	t.Environment["BASEASSET"] = t.Config.BaseAsset
	t.Environment["QUOTEASSET"] = t.Config.QuoteAsset
	t.Environment["CEX_BASEASSET"] = t.Config.CexBaseAsset
	t.Environment["CEX_QUOTEASSET"] = t.Config.CexQuoteAsset
	t.Environment["OPENDEX_RPC_PORT"] = rpcPort
	t.Environment["CEX"] = t.Config.Cex
	t.Environment["CEX_API_SECRET"] = t.Config.CexApiSecret
	t.Environment["CEX_API_KEY"] = t.Config.CexApiKey
	t.Environment["TEST_MODE"] = fmt.Sprint(t.Config.TestMode)
	t.Environment["MARGIN"] = t.Config.Margin
	t.Environment["TEST_CENTRALIZED_EXCHANGE_BASEASSET_BALANCE"] = t.Config.TestCentralizedBaseassetBalance
	t.Environment["TEST_CENTRALIZED_EXCHANGE_QUOTEASSET_BALANCE"] = t.Config.TestCentralizedQuoteassetBalance

	return nil
}
