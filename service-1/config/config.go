package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
	Nats     Nats     `yaml:"nats"`
}

type Nats struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Database struct {
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
}

type Postgres struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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
