package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	GRPC struct {
		Host string
		Port int
	}
}

func NewAppConfig(configFile string) (*AppConfig, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfig(configFile string) (*AppConfig, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, err
}
