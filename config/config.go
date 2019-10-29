package config

import (
	"github.com/blocktop/mp-common/config"
	"github.com/caarlos0/env"
)

type Config struct {
	config.BaseConfig
	WebAuthEndpoint          string `env:"MP_WEB_AUTH_ENDPOINT"`
	AnchorName               string `env:"MP_ANCHOR_NAME" envDefault:"Blocktop"`
}

var cfg *Config

func init() {
	cfg = &Config{}
	cfg.Parse()
}

func GetConfig() *Config {
	return cfg
}

func (c *Config) Parse() {
	c.BaseConfig.Parse()

	err := env.Parse(c)
	if err != nil {
		panic(err)
	}
}