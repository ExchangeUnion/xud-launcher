package service

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

type GethConfig struct {
	Mode                string `usage:"Geth service mode"`
	Rpchost             string `usage:"External geth RPC hostname"`
	Rpcport             uint16 `usage:"External geth RPC port"`
	InfuraProjectId     string `usage:"Infura geth provider project ID"`
	InfuraProjectSecret string `usage:"Infura geth provider project secret"`
	Cache               string `usage:"Geth cache size"`
	AncientChaindataDir string `usage:"Specify the container's volume mapping ancient chaindata directory. Can be located on a slower HDD."`
}

type Geth struct {
	Base

	config GethConfig
}

func newGeth(name string) Geth {
	return Geth{
		Base: newBase(name),
	}
}

func (t *Geth) ConfigureFlags(cmd *cobra.Command) error {
	err := t.Base.ConfigureFlags(&BaseConfig{
		Disable:     false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       "exchangeunion/geth",
	}, cmd)
	if err != nil {
		return err
	}

	if err := ReflectFlags(t.Name, &t.config, cmd); err != nil {
		return err
	}

	//cmd.PersistentFlags().StringVar(&t.config.Mode, "geth.mode", "light", "Geth service mode")
	//cmd.PersistentFlags().StringVar(&t.config.Rpchost, "geth.rpchost", "", "External geth RPC hostname")
	//cmd.PersistentFlags().Uint16Var(&t.config.Rpcport, "geth.rpcport", 0, "External geth RPC port")
	//cmd.PersistentFlags().StringVar(&t.config.InfuraProjectId, "geth.infura-project-id", "", "Infura geth provider project ID")
	//cmd.PersistentFlags().StringVar(&t.config.InfuraProjectSecret, "geth.infura-project-secret", "", "Infura geth provider project secret")
	//cmd.PersistentFlags().StringVar(&t.config.Cache, "geth.cache", "", "Geth cache size")
	//cmd.PersistentFlags().StringVar(&t.config.AncientChaindataDir, "geth.ancient-chaindata-dir", "", "Specify the container's volume mapping ancient chaindata directory. Can be located on a slower HDD.")

	return nil
}

func (t *Geth) GetConfig() interface{} {
	return t.config
}

func (t *Geth) Apply(config *SharedConfig, services map[string]Service) error {
	network := config.Network

	// validation
	if network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/root/.ethereum", network)
	if err != nil {
		return err
	}

	// geth apply
	t.Environment["NETWORK"] = network

	if t.config.AncientChaindataDir != "" {
		volume := fmt.Sprintf("%s:/root/.ethereum-ancient-chaindata", t.config.AncientChaindataDir)
		t.Volumes = append(t.Volumes, volume)
		t.Environment["CUSTOM_ANCIENT_CHAINDATA"] = "true"
	}

	if t.config.Cache != "" {
		t.Command = append(t.Command, fmt.Sprintf("--cache %s", t.config.Cache))
	}

	// TODO select ethProvider in light mode

	return nil
}
