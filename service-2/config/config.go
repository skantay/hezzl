package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database Database `yaml:"database"`
}

type Database struct {
	Clickhouse Clickhouse `yaml:"clickhouse"`
}

type Clickhouse struct {
	User     string `yaml:"user"`
	Password string `yaml:"default_password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

func Load(path string) (Config, error) {
	viper.SetConfigFile(path)

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
