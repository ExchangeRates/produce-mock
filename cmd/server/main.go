package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/ExchangeRates/produce-mock/internal"
	"github.com/ExchangeRates/produce-mock/internal/config"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/service.toml", "path to config file")
}

func main() {
	flag.Parse()

	configuration := config.NewConfig()

	if _, err := toml.DecodeFile(configPath, configuration); err != nil {
		log.Fatalln(err)
	}

	if err := internal.Start(configuration); err != nil {
		log.Fatalln(err)
	}
}
