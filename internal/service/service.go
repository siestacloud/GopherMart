package service

import "github.com/siestacloud/gopherMart/internal/repository"

// Главный тип слоя SVC, который встраивается как зависимость в слое TRANSPORT
type Service struct {
}

// Конструктор слоя SVC
func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
