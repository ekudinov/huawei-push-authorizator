package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ApiToken        string `env:"HU_API_TOKEN"`
	ClientID        string `env:"HU_CLIENT_ID"`
	ClientSecret    string `env:"HU_CLIENT_SECRET"`
	Host            string `env:"HU_HOST" envDefault:"localhost"`
	Port            int    `env:"HU_PORT" envDefault:"8077"`
	CheckInterval   int    `env:"HU_CHECK_INTERVAL" envDefault:"15"`
	EarlyUpdateTime int    `env:"HU_EARLY_UPDATE_TIME" envDefault:"30"`
}

func NewConfigFromEnv() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Can not create app configuration from env!", err)
	}

	log.Println("Load config ok.")

	return &cfg
}
