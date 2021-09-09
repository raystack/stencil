package config

import "github.com/odpf/salt/config"

func Load(configFile string) (Config, error) {
	var cfg Config
	loader := config.NewLoader(config.WithFile(configFile))

	if err := loader.Load(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
