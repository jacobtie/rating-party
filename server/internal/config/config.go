package config

import (
	"fmt"
	"time"

	"github.com/jacobtie/rating-party/server/internal/platform/namegenerator"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Instance    string
	Environment ENV `default:"local" envconfig:"APP_ENV"`
	Web         struct {
		APIHost      string        `default:":3000" envconfig:"API_HOST"`
		ReadTimeout  time.Duration `default:"10s" envconfig:"READ_TIMEOUT"`
		WriteTimeout time.Duration `default:"15s" envconfig:"WRITE_TIMEOUT"`
	}
	DB struct {
		DriverName string `default:"mysql" envconfig:"DRIVER_NAME"`
		DBUser     string `default:"root" envconfig:"DB_USER"`
		DBPass     string `default:"" envconfig:"DB_PASS"`
		DBURI      string `default:"localhost:3306" envconfig:"DB_URI"`
		DBName     string `default:"ratingparty" envconfig:"DB_NAME"`
	}
	AdminPasscode  string `default:"ivory" envconfig:"ADMIN_PASSCODE"`
	AdminJWTSecret string `default:"ebony" envconfig:"ADMIN_JWT_SECRET"`
}

type ENV string

const (
	ENV_LOCAL ENV = "local"
	ENV_PROD  ENV = "production"
)

var (
	cfg *Config
)

func Get() (*Config, error) {
	if cfg == nil {
		var config Config
		if err := envconfig.Process("", &config); err != nil {
			return nil, fmt.Errorf("[db.Get] failed to process env vars: %w", err)
		}
		cfg = &config
		cfg.Instance = namegenerator.Generator()
	}
	return cfg, nil
}
