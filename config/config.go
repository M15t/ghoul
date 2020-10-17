package config

import (
	"fmt"
	"os"

	cfgutil "github.com/M15t/ghoul/pkg/util/cfg"
)

// Configuration holds data necessery for configuring application
type Configuration struct {
	Stage        string   `env:"UP_STAGE"`
	Host         string   `env:"HOST"`
	Port         int      `env:"PORT"`
	ReadTimeout  int      `env:"READ_TIMEOUT"`
	WriteTimeout int      `env:"WRITE_TIMEOUT"`
	AllowOrigins []string `env:"ALLOW_ORIGINS"`
	Debug        bool     `env:"DEBUG"`
	DbLog        bool     `env:"DB_LOG"`
	DbPsn        string   `env:"DB_PSN"`
	JwtSecret    string   `env:"JWT_SECRET"`
	JwtDuration  int      `env:"JWT_DURATION"`
	JwtAlgorithm string   `env:"JWT_ALGORITHM"`
}

// Load returns Configuration struct
func Load() (*Configuration, error) {
	appName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if configname := os.Getenv("CONFIG_NAME"); configname != "" {
		appName = configname
	}
	stage := os.Getenv("UP_STAGE")
	if configstage := os.Getenv("CONFIG_STAGE"); configstage != "" {
		stage = configstage
	}

	cfg := new(Configuration)
	if err := cfgutil.LoadWithAPS(cfg, appName, stage); err != nil {
		return nil, fmt.Errorf("Error parsing environment config: %s", err)
	}
	return cfg, nil
}
