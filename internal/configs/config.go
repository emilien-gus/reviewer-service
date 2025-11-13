package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Token struct {
	Token string `mapstructure:"token"`
	Role  string `mapstructure:"role"`
}

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`
	Auth struct {
		Tokens []Token `mapstructure:"tokens"`
	} `mapstructure:"auth"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}
