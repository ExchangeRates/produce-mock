package config

import "github.com/sirupsen/logrus"

type Config struct {
	Port     int          `toml:"port"`
	LogLevel logrus.Level `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		Port:     8080,
		LogLevel: logrus.DebugLevel,
	}
}
