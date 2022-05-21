package config

import (
	"github.com/spf13/viper"
)

// InitWithFile init config with config file
func InitWithFile(filename string) error {
	viper.SetConfigFile(filename)

	return viper.ReadInConfig()
}
