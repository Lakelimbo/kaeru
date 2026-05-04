package config

import (
	"os"

	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/goccy/go-yaml"
)

// Load gets the app's basic configuration, such as port, SQLite path,
// among others.
//
// If no config is found, or specific, non-nil fields are not set, it will
// fallback to some default values.
func Load(path string) (*Config, error) {
	c := defaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warnf("Local config file not found. Using some default settings...")
			return &c, nil
		}

		return nil, err
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	logger.Infof("Configuration loaded successfully (from %s)", path)
	return &c, nil
}

func defaultConfig() Config {
	env := os.Getenv(ENVIRONMENT)
	if env == "" {
		env = "development"
	}

	return Config{
		Name:        "kaeru instance",
		Environment: env,
		OutputDir:   ".kaeru",
		Server: ServerConfig{
			Host:              "127.0.0.1",
			Port:              4040,
			ComposePath:       "./.kaeru/compose/",
			APISpecProduction: false,
		},
		Database: DatabaseConfig{
			Path: "./.kaeru/db/kaeru.db",
			WAL:  true,
		},
		Logging: LoggingConfig{
			Level: "debug",
		},
	}
}
