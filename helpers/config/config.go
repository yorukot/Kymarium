package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

type AppEnv string

const (
	AppEnvDev  AppEnv = "dev"
	AppEnvProd AppEnv = "prod"
)

// EnvConfig holds all environment variables for the application
type EnvConfig struct {
	// PostgreSQL Settings
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSL_MODE,required"`

	// Redis/Dragonfly Settings
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`

	AppEnv       AppEnv `env:"APP_ENV" envDefault:"prod"`
	AppName      string `env:"APP_NAME" envDefault:"knocker"`
	AppMachineID int16    `env:"APP_MACHINE_ID" envDefault:"0"`
}

var (
	appConfig *EnvConfig
	once      sync.Once
)

// loadConfig loads and validates all environment variables
func loadConfig() (*EnvConfig, error) {
	cfg := &EnvConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// InitConfig initializes the config only once
func InitConfig() (*EnvConfig, error) {
	var err error
	once.Do(func() {
		appConfig, err = loadConfig()
		zap.L().Info("Config loaded")
	})
	return appConfig, err
}

// Env returns the config. Panics if not initialized.
func Env() *EnvConfig {
	if appConfig == nil {
		zap.L().Panic("config not initialized — call InitConfig() first")
	}
	return appConfig
}
