package util

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Config struct {
	Host  string
	Pics  string
	Mongo ConfigMongo
}

type ConfigMongo struct {
	Uri       string
	Db        string
	Collction string
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(cfg.Host, "/") {
		cfg.Host = cfg.Host[:len(cfg.Host)-1]
	}

	return cfg, nil
}
