package handler

import (
	"github.com/siestacloud/gopherMart/internal/config"
	"github.com/siestacloud/gopherMart/internal/service"
)

type Handler struct {
	cfg      *config.Cfg
	services *service.Service
}

func NewHandler(cfg *config.Cfg, services *service.Service) *Handler {
	return &Handler{
		cfg:      cfg,
		services: services,
	}
}
