// Package config provides configuration settings
// mapped to go structs.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppCfg struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`

	JWTExpSecs int64 `mapstructure:"jwt_exp_secs"`
}

func LoadConfig(relPath string) AppCfg {
	if relPath == "" {
		relPath = "app_config.json"
	}

	viper.SetConfigFile(relPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Error loading config file: %w \n", err))
	}

	var cfg AppCfg

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("Unable to decode into struct: %w", err))
	}

	return cfg
}
