package config

import (
	"github.com/spf13/viper"
)

func LoadConfig(config *Config, args *Arguments) error {
	viper.SetConfigFile(args.ConfigPath)
	viper.SetConfigType("toml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&config)

	return err

}
