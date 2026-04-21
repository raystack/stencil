package config

import (
	"github.com/raystack/salt/config"
)

// Load reads the configuration from the given file path and returns a Config.
func Load(configFile string) (Config, error) {
	var cfg Config
	loader := config.NewLoader(config.WithFile(configFile))

	if err := loader.Load(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
