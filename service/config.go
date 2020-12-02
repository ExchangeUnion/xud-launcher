package service

var (
	images = map[string]map[string]string{
		"simnet": {
			"lndbtc":  "exchangeunion/lndbtc-simnet:latest",
			"lndltc":  "exchangeunion/lndltc-simnet:latest",
			"connext": "connextproject/vector_node:837bafa1",
			"arby":    "exchangeunion/arby:latest",
			"webui":   "exchangeunion/webui:latest",
			"proxy":   "exchangeunion/proxy:latest__proxy-first",
			"xud":     "exchangeunion/xud:latest",
		},
		"testnet": {
			"bitcoind":  "exchangeunion/bitcoind:latest",
			"litecoind": "exchangeunion/litecoind:latest",
			"geth":      "exchangeunion/geth:latest",
			"lndbtc":    "exchangeunion/lndbtc:latest",
			"lndltc":    "exchangeunion/lndltc:latest",
			"connext":   "exchangeunion/connext:latest",
			"arby":      "exchangeunion/arby:latest",
			"boltz":     "exchangeunion/boltz:latest",
			"webui":     "exchangeunion/webui:latest",
			"proxy":     "exchangeunion/proxy:latest__proxy-first",
			"xud":       "exchangeunion/xud:latest",
		},
		"mainnet": {
			"bitcoind":  "exchangeunion/bitcoind:0.20.1",
			"litecoind": "exchangeunion/litecoind:0.18.1",
			"geth":      "exchangeunion/geth:1.9.24",
			"lndbtc":    "exchangeunion/lndbtc:0.11.1-beta",
			"lndltc":    "exchangeunion/lndltc:0.11.0-beta.rc1",
			"connext":   "exchangeunion/connext:1.3.6",
			"arby":      "exchangeunion/arby:1.3.0",
			"boltz":     "exchangeunion/boltz:1.1.1",
			"webui":     "exchangeunion/webui:1.0.0",
			"proxy":     "exchangeunion/proxy:latest__proxy-first",
			"xud":       "exchangeunion/xud:1.2.0",
		},
	}

	ethProviders = map[string][]string{
		"testnet": {
			"http://eth.kilrau.com:52041",
			"http://michael1011.at:8546",
			"http://gethxudxv2k4pv5t5a5lswq2hcv3icmj3uwg7m2n2vuykiyv77legiad.onion:8546",
		},
		"mainnet": {
			"http://eth.kilrau.com:41007",
			"http://michael1011.at:8545",
			"http://gethxudxv2k4pv5t5a5lswq2hcv3icmj3uwg7m2n2vuykiyv77legiad.onion:8545",
		},
	}
)
