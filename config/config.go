package config

import (
	"github.com/naoina/toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type Config struct {
	DbHost     string `toml:"db_host"`
	DbSpace    string `toml:"db_space"`
	DbUser     string `toml:"db_username"`
	DbPasswd   string `toml:"db_password"`
	GpcAddress string `toml:"grpc_addr"`
}

var _cfg *Config = nil

func ParseConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("get config failed", "err", err)
		panic(err)
	}
	err = toml.Unmarshal(data, &_cfg)
	if err != nil {
		log.Error("unmarshal config failed", "err", err)
		panic(err)
	}
	return _cfg, nil
}

func GetConfig() *Config {
	return _cfg
}
