package config

import "github.com/sirupsen/logrus"

type Config struct {
	Port     int          `toml:"port"`
	LogLevel logrus.Level `toml:"log_level"`
	Urls     []string     `toml:"urls"`
}

func NewConfig() *Config {
	return &Config{
		Port:     8080,
		LogLevel: logrus.DebugLevel,
		Urls:     []string{"localhost:9092"},
	}
}
