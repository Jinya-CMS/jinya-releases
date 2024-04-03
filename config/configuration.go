package config

import (
	"fmt"
	"go-simpler.org/env"
)

type Configuration struct {
	PostgresUrl      string `env:"DATABASE_URL"`
	OpenIDClientId   string `env:"OPENID_CLIENT_ID"`
	OpenIDDomain     string `env:"OPENID_DOMAIN"`
	ServerUrl        string `env:"SERVER_URL"`
	StorageSecretKey string `env:"STORAGE_SECRET_KEY"`
	StorageAccessKey string `env:"STORAGE_ACCESS_KEY"`
	StorageBucket    string `env:"STORAGE_BUCKET"`
	StorageUrl       string `env:"STORAGE_URL"`
}

func (c Configuration) GetRedirectUrl() string {
	return fmt.Sprintf("%s/admin/login/callback", c.ServerUrl)
}

var LoadedConfiguration *Configuration

func LoadConfiguration() error {
	config := new(Configuration)
	err := env.Load(config, nil)
	if err != nil {
		return err
	}

	LoadedConfiguration = config

	return nil
}
