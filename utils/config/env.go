package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload" // Load .env into the process environment.
	"github.com/yorukot/kymarium/models"
)

// AppEnv represents the running environment.
type AppEnv string

// AppEnv values.
const (
	AppEnvDev  AppEnv = "dev"
	AppEnvProd AppEnv = "prod"
)

// EnvConfig holds all environment variables for the application
type EnvConfig struct {
	AppEnv       AppEnv   `env:"APP_ENV" envDefault:"prod"`
	AppName      string   `env:"APP_NAME" envDefault:"kymarium"`
	AppMachineID int16    `env:"APP_MACHINE_ID" envDefault:"1"`
	AppPort      string   `env:"APP_PORT" envDefault:"8000"`
	AppRegion    string   `env:"APP_REGION" envDefault:"TW-Taipei"`
	AppRegions   []string `env:"APP_REGIONS" envDefault:"TW-Taipei" envSeparator:","`

	// Security Settings
	JWTSecretKey   string `env:"JWT_SECRET_KEY,required" envDefault:"change_me_to_a_secure_key"`
	FrontendDomain string `env:"FRONTEND_DOMAIN" envDefault:"localhost"`
	CookieDomain   string `env:"COOKIE_DOMAIN" envDefault:""`

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

	// Optional Settings
	OAuthStateExpiresAt   int `env:"OAUTH_STATE_EXPIRES_AT" envDefault:"600"`        // 10 minutes
	AccessTokenExpiresAt  int `env:"ACCESS_TOKEN_EXPIRES_AT" envDefault:"900"`       // 15 minutes
	RefreshTokenExpiresAt int `env:"REFRESH_TOKEN_EXPIRES_AT" envDefault:"31536000"` // 365 days
	SessionExpiresAt      int `env:"SESSION_EXPIRES_AT" envDefault:"432000"`         // 5 days

	// Email Verification
	EmailVerifyExpiresAt int    `env:"EMAIL_VERIFY_EXPIRES_AT" envDefault:"900"` // 15 minutes
	BackendURL           string `env:"BACKEND_URL" envDefault:"http://localhost:8000"`

	// SMTP Settings
	SMTPEnabled  bool   `env:"SMTP_ENABLED" envDefault:"false"`
	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     string `env:"SMTP_PORT"`
	SMTPUsername string `env:"SMTP_USERNAME"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SMTPFrom     string `env:"SMTP_FROM"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID,required"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET,required"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL,required"`
}

var (
	appConfig    *EnvConfig
	once         sync.Once
	regionsOnce  sync.Once
	regionsErr   error
	regionsByID  map[int64]models.Region
	regionsByKey map[string]models.Region
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
	})
	return appConfig, err
}

// Env returns the config. Panics if not initialized.
func Env() *EnvConfig {
	if appConfig == nil {
		panic("config not initialized — call InitConfig() first")
	}
	return appConfig
}
