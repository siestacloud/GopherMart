package main

import (
	"log"

	"github.com/siestacloud/gopherMart/internal/config"
	"github.com/siestacloud/gopherMart/internal/repository"
	"github.com/siestacloud/gopherMart/internal/service"
	"github.com/siestacloud/gopherMart/internal/transport/rest"
	"github.com/siestacloud/gopherMart/internal/transport/rest/handler"
	"github.com/sirupsen/logrus"
)

var (
	cfg config.Cfg
)

// @title Template App API
// @version 1.0
// @description API Server for Template Application

// @host localhost:9999
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	err := config.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgresDB(cfg.URLPostgres)
	if err != nil {
		logrus.Warnf("failed to initialize postrges: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(&cfg, services)

	s, err := rest.NewServer(&cfg, handlers)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

}
