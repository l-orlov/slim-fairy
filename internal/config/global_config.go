package config

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

// TODO: use private config without write access

type (
	Config struct {
		PgDSN         string `mapstructure:"PG_DSN"`
		ServerAddress string `mapstructure:"SERVER_ADDRESS"`
		JWTToken      `mapstructure:",squash"`
		TelegramBot   `mapstructure:",squash"`
		AIConfig      `mapstructure:",squash"`
	}
	JWTToken struct {
		Lifitime time.Duration `mapstructure:"TOKEN_LIFETIME"`
		Issuer   string        `mapstructure:"TOKEN_ISSUER"`
		Secret   string        `mapstructure:"TOKEN_SECRET"`
	}
	TelegramBot struct {
		Token string `mapstructure:"TG_BOT_TOKEN"`
	}
	AIConfig struct {
		APIKey string `mapstructure:"API_KEY"`
	}
)

var (
	once   sync.Once
	config *Config
)

// Load reads configuration from file or environment variables.
func Load(path string) (err error) {
	once.Do(func() {
		viper.AddConfigPath(path)
		viper.SetConfigName("app_prod")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		err = viper.Unmarshal(&config)
	})

	return
}

// Get returns global config
func Get() *Config {
	return config
}
