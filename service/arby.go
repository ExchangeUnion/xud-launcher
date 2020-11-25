package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type ArbyConfig struct {
	BaseConfig

	LiveCex                          bool
	TestMode                         bool
	BaseAsset                        string
	QuoteAsset                       string
	CexBaseAsset                     string
	CexQuoteAsset                    string
	TestCentralizedBaseassetBalance  string
	TestCentralizedQuoteassetBalance string
	Cex                              string
	CexApiKey                        string
	CexApiSecret                     string
	Margin                           string
}

type Arby struct {
	Base

	config ArbyConfig
}

func NewArby() Arby {
	return Arby{
		config: ArbyConfig{},
	}
}

func (t Arby) ConfigureFlags(cmd *cobra.Command) error {
	err := configureCommonFlags("arby", &t.config.BaseConfig, cmd)
	if err != nil {
		return err
	}

	cmd.PersistentFlags().BoolVar(&t.config.LiveCex, "arby.live-cex", true, "Live CEX (deprecated)")
	err = cmd.PersistentFlags().MarkDeprecated("arby.live-cex", "Please use --arby.test-mode instead")
	if err != nil {
		return err
	}
	cmd.PersistentFlags().BoolVar(&t.config.TestMode, "arby.test-mode", true, "Whether to issue real orders on the centralized exchange")
	cmd.PersistentFlags().StringVar(&t.config.BaseAsset, "arby.base-asset", "", "Base asset")
	cmd.PersistentFlags().StringVar(&t.config.QuoteAsset, "arby.quote-asset", "", "Quote asset")
	cmd.PersistentFlags().StringVar(&t.config.CexBaseAsset, "arby.cex-base-asset", "", "Centralized exchange base asset")
	cmd.PersistentFlags().StringVar(&t.config.CexQuoteAsset, "arby.cex-quote-asset", "", "Centralized exchange quote asset")
	cmd.PersistentFlags().StringVar(&t.config.TestCentralizedBaseassetBalance, "arby.test-centralized-baseasset-balance", "", "Test centralized base asset balance")
	cmd.PersistentFlags().StringVar(&t.config.TestCentralizedQuoteassetBalance, "arby.test-centralized-quoteasset-balance", "", "Test centralized quote asset balance")
	cmd.PersistentFlags().StringVar(&t.config.Cex, "arby.cex", "binance", "Centralized Exchange")
	cmd.PersistentFlags().StringVar(&t.config.CexApiKey, "arby.cex-api-key", "123", "CEX API key")
	cmd.PersistentFlags().StringVar(&t.config.CexApiSecret, "arby.cex-api-secret", "abc", "CEX API secret")
	cmd.PersistentFlags().StringVar(&t.config.Margin, "arby.margin", "0.04", "Trade margin")

	return nil
}

func (t Arby) Apply(network string, services map[string]Service) error {

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
	t.Environment["BASEASSET"] = t.config.BaseAsset
	t.Environment["QUOTEASSET"] = t.config.QuoteAsset
	t.Environment["CEX_BASEASSET"] = t.config.CexBaseAsset
	t.Environment["CEX_QUOTEASSET"] = t.config.CexQuoteAsset
	t.Environment["OPENDEX_RPC_PORT"] = rpcPort
	t.Environment["CEX"] = t.config.Cex
	t.Environment["CEX_API_SECRET"] = t.config.CexApiSecret
	t.Environment["CEX_API_KEY"] = t.config.CexApiKey
	t.Environment["TEST_MODE"] = fmt.Sprint(t.config.TestMode)
	t.Environment["MARGIN"] = t.config.Margin
	t.Environment["TEST_CENTRALIZED_EXCHANGE_BASEASSET_BALANCE"] = t.config.TestCentralizedBaseassetBalance
	t.Environment["TEST_CENTRALIZED_EXCHANGE_QUOTEASSET_BALANCE"] = t.config.TestCentralizedQuoteassetBalance

	return nil
}
