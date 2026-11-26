package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/itsviktor/sir/internal/utils"
)

type DriverType string

const (
	DriverSqlite DriverType = "sqlite"
)

type Config struct {
	Driver DriverType `yaml:"driver"`
}

func validateDriver(cfg *Config) error {
	switch cfg.Driver {
	case "sqlite":
		return nil
	default:
		return fmt.Errorf("invalid driver value: %s", cfg.Driver)
	}
}

func Load() *Config {
	content, err := os.ReadFile("sir.yaml")
	if os.IsNotExist(err) {
		utils.Fatal("config file sir.yaml is required")
	} else if err != nil {
		utils.Fatal("error while reading config file sir.yaml", "description", err.Error())
	}

	c := new(Config)
	if err := yaml.Unmarshal(content, c); err != nil {
		utils.Fatal("failed to load config file sir.yaml", "description", err.Error())
	}

	if err := validateDriver(c); err != nil {
		utils.Fatal("config validation error", "description", err.Error())
	}

	return c
}
