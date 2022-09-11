package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

// Setting up configuration parametrs
func Parse(cfg *Cfg) error {

	// Читаю флаги, переопределяю параметры, если флаги заданы

	flag.StringVar(&cfg.Address, "a", "localhost:8080", "Address for server. Possible values: localhost:8080")
	flag.StringVar(&cfg.URLPostgres, "d", "not set", "url for postgres db con. Possible values: url")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "http://localhost:8081", "black api. Possible values: http://localhost:8081")
	flag.StringVar(&cfg.Logrus.LogLevel, "l", "info", "log level. Possible values: debug, info")
	flag.BoolVar(&cfg.Logrus.JSON, "f", false, "JSON log format. Possible values: true, false")

	flag.Parse()

	// Читаю переменные окружения, переопределяю параметры, если пер окр заданы
	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
