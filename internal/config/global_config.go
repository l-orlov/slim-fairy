package config

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		PgDSN         string `mapstructure:"PG_DSN"`
		ServerAddress string `mapstructure:"SERVER_ADDRESS"`
		JWTToken      `mapstructure:",squash"`
	}
	JWTToken struct {
		Lifitime time.Duration `mapstructure:"TOKEN_LIFETIME"`
		Issuer   string        `mapstructure:"TOKEN_ISSUER"`
		Secret   string        `mapstructure:"TOKEN_SECRET"`
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
		viper.SetConfigName("app")
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
