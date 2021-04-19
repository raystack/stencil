package config

import (
	"fmt"
	"strings"

	"github.com/jeremywohl/flatten"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

//Config Server config
type Config struct {
	Port            string `default:"8080"`
	Timeout         int    `default:"60"`
	BucketURL       string
	Scheme          string `default:"https"`
	Host            string `default:"localhost"`
	AuthBearerToken string
}

// LoadConfig returns application configuration
func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file. Envs will be used")
		} else {
			panic(err)
		}
	}

	err = bindEnvKeys(Config{})
	throwError(err, "Unable to bind env keys")

	var config Config
	defaults.SetDefaults(&config)
	err = viper.Unmarshal(&config)
	throwError(err, "viper unmarshal failed")
	return &config
}

// viper.Unmarshal doesn't work directly on envs. We have to bind them manually. See https://github.com/spf13/viper/issues/584
func bindEnvKeys(config Config) error {
	var structMap map[string]interface{}
	err := mapstructure.Decode(config, &structMap)
	if err != nil {
		return err
	}

	flat, err := flatten.Flatten(structMap, "", flatten.DotStyle)
	if err != nil {
		return err
	}

	for key := range flat {
		viper.BindEnv(key)
	}
	return nil
}

func throwError(err error, message string) {
	if err != nil {
		err = fmt.Errorf("%s\n%v", message, err)
		panic(err)
	}
}
