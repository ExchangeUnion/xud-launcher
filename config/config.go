package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

type GitHub struct {
	AccessToken string
}

type Config struct {
	GitHub GitHub
	SimnetDir string
	TestnetDir string
	MainnetDir string
}

func ParseConfig(configFile string) (*Config, error) {
	config := Config{}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &config, nil
}
