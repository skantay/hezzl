package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
}

type Database struct {
	Postgres   Postgres   `yaml:"postgres"`
	Redis      Redis      `yaml:"redis"`
	Clickhouse Clickhouse `yaml:"clickhouse"`
}

type Clickhouse struct {
	User     string `yaml:"user"`
	Password string `yaml:"default_password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Sslmode  string `yaml:"sslmode"`
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
