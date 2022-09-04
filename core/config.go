package core

import (
	"log"

	env "github.com/Netflix/go-env"
)

type Config struct {
	Auth0 struct {
		AUTH0_DOMAIN        string `env:"AUTH0_DOMAIN"`
		AUTH0_CLIENT_ID     string `env:"AUTH0_CLIENT_ID"`
		AUTH0_CLIENT_SECRET string `env:"AUTH0_CLIENT_SECRET"`
		AUTH0_CALLBACK_URL  string `env:"AUTH0_CALLBACK_URL"`
	}

	Host struct {
		PORT string `env:"PORT"`
	}

	Extras env.EnvSet
}

func GetConfig() *Config {
	var conf Config
	es, err := env.UnmarshalFromEnviron(&conf)
	if err != nil {
		log.Fatal(err)
	}
	// Remaining environment variables.
	conf.Extras = es
	return &conf
}
