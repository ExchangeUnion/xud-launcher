package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
)

type ConnextConfig struct {
	// add more connext specified attributes here
}

type Connext struct {
	Base

	Config ConnextConfig
}

func newConnext(name string) Connext {
	return Connext{
		Base: newBase(name),
	}
}

func (t *Connext) ConfigureFlags(cmd *cobra.Command, network string) error {
	if err := t.Base.ConfigureFlags(cmd, network, &BaseConfig{
		Disabled:    false,
		ExposePorts: []string{},
		Dir:         fmt.Sprintf("./data/%s", t.Name),
		Image:       images[network][t.Name],
	}); err != nil {
		return err
	}

	// configure connext specified flags here

	return nil
}

func (t *Connext) GetConfig() interface{} {
	return t.Config
}

func (t *Connext) Apply(config *SharedConfig, services map[string]Service) error {
	ReflectFillConfig(t.Name, &t.Config)

	network := config.Network

	// validation
	if network != "simnet" && network != "testnet" && network != "mainnet" {
		return errors.New("invalid network: " + network)
	}

	// base apply
	err := t.Base.Apply("/app/connext-store", config.Network)
	if err != nil {
		return err
	}

	t.Environment["NETWORK"] = network

	// connext apply
	if strings.Contains(t.Image, "vector_node") {
		t.Environment["VECTOR_CONFIG"] = `{
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
}`
		t.Environment["VECTOR_SQLITE_FILE"] = "/database/store.db"
		t.Environment["VECTOR_PROD"] = "true"
	} else {
		t.Environment["LEGACY_MODE"] = "true"
		switch network {
		case "simnet":
			t.Environment["CONNEXT_ETH_PROVIDER_URL"] = "http://connext.simnet.exchangeunion.com:8545"
			t.Environment["CONNEXT_NODE_URL"] = "https://connext.simnet.exchangeunion.com"
		case "testnet":
			t.Environment["CONNEXT_NODE_URL"] = "https://connext.testnet.exchangeunion.com"
		case "mainnet":
			t.Environment["CONNEXT_NODE_URL"] = "https://connext.boltz.exchange"
		}

		gethConfig := services["geth"].GetConfig().(GethConfig)

		mode := gethConfig.Mode
		switch mode {
		case "external":
			rpcHost := gethConfig.Rpchost
			rpcPort := gethConfig.Rpcport
			t.Environment["CONNEXT_ETH_PROVIDER_URL"] = fmt.Sprintf("http://%s:%d", rpcHost, rpcPort)
		case "infura":
			projId := gethConfig.InfuraProjectId
			switch network {
			case "mainnet":
				t.Environment["CONNEXT_ETH_PROVIDER_URL"] = fmt.Sprintf("https://mainnet.infura.io/v3/%s", projId)
			case "testnet":
				t.Environment["CONNEXT_ETH_PROVIDER_URL"] = fmt.Sprintf("https://rinkeby.infura.io/v3/%s", projId)
			case "simnet":
				return errors.New("no Infura Ethereum provider for simnet")
			}
		case "light":
			t.Environment["CONNEXT_ETH_PROVIDER_URL"] = selectFastestProvider(ethProviders[network])
		case "native":
			t.Environment["CONNEXT_ETH_PROVIDER_URL"] = "http://geth:8545"
		}
	}

	return nil
}

func checkProvider(url string) error {
	var payload = []byte(`{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}`)
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	return nil
}

func selectFastestProvider(providers []string) string {
	return providers[0]
}

func (t *Connext) ToJson() map[string]interface{} {
	result := t.Base.ToJson()

	rpc := make(map[string]interface{})
	result["rpc"] = rpc
	rpc["type"] = "HTTP"
	rpc["host"] = "connext"
	rpc["port"] = 5040

	return result
}
