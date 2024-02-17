package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/FreshOfficeFriends/SSO/pkg/database"
	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

const (
	pathToEnv = "../../.env"
)

func init() {
	err := godotenv.Load(pathToEnv)
	if err != nil {
		os.Exit(1)
	}
}

type Config struct {
	DB database.Config
}

var (
	cfg  = new(Config)
	once sync.Once
)

func New() *Config {
	once.Do(func() {
		if err := envconfig.Process("DB", &cfg.DB); err != nil {
			logger.Fatal(err.Error())
		}
	})
	return cfg
}
