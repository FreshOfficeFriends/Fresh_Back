package config

import (
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/FreshOfficeFriends/SSO/pkg/database"
	"github.com/FreshOfficeFriends/SSO/pkg/logger"
)

//todo избавиться от os.getenv в коде, все должно быть здесь.. Если что прокидывай в контекст

func init() {
	err := godotenv.Load()
	if err != nil {
		os.Exit(15)
	}
}

type Config struct {
	DB  database.Config
	JWT JWTConfig
}

type JWTConfig struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Secret     []byte
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

		if err := envconfig.Process("JWT", &cfg.JWT); err != nil {
			logger.Fatal(err.Error())
		}
	})
	return cfg
}
