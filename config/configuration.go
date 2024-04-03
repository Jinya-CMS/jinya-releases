package config

import (
	"fmt"
	"go-simpler.org/env"
)

type Configuration struct {
	MongoUrl       string `env:"DATABASE_URL"`
	OpenIDClientId string `env:"OPENID_CLIENT_ID"`
	OpenIDDomain   string `env:"OPENID_DOMAIN"`
	ServerUrl      string `env:"SERVER_URL"`
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
