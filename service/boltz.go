package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type BoltzConfig struct {
	// add more boltz specified attributes here
}

type Boltz struct {
	Base

	Config BoltzConfig
}

func newBoltz(name string) Boltz {
	return Boltz{
		Base: newBase(name),
	}
}

func (t *Boltz) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network, &BaseConfig{
		Disabled:    network == "simnet",
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       images[network][t.Name],
	}); err != nil {
		return err
	}

	// configure boltz specified flags here

	return nil
}

func (t *Boltz) GetConfig() interface{} {
	return t.Config
}

func (t *Boltz) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.boltz", config.Network)
	if err != nil {
		return err
	}

	// boltz apply

	t.Volumes = append(t.Volumes,
		"./data/lndbtc:/root/.lndbtc",
		"./data/lndltc:/root/.lndltc",
	)

	return nil
}

func (t *Boltz) ToJson() map[string]interface{} {
	result := t.Base.ToJson()

	rpc := make(map[string]interface{})
	result["rpc"] = rpc

	bitcoin := make(map[string]interface{})
	bitcoin["type"] = "gRPC"
	bitcoin["host"] = "boltz"
	bitcoin["port"] = 9002
	bitcoin["tlsCert"] = fmt.Sprintf("%s/%s/bitcoin/tls.cert", PROXY_DATA_DIR, t.GetName())
	bitcoin["macaroon"] = fmt.Sprintf("%s/%s/bitcoin/admin.macaroon", PROXY_DATA_DIR, t.GetName())

	litecoin := make(map[string]interface{})
	litecoin["type"] = "gRPC"
	litecoin["host"] = "boltz"
	litecoin["port"] = 9102
	litecoin["tlsCert"] = fmt.Sprintf("%s/%s/litecoin/tls.cert", PROXY_DATA_DIR, t.GetName())
	litecoin["macaroon"] = fmt.Sprintf("%s/%s/litecoin/admin.macaroon", PROXY_DATA_DIR, t.GetName())

	rpc["bitcoin"] = bitcoin
	rpc["litecoin"] = litecoin

	return result
}
