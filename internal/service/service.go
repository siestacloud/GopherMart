package service

import (
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/internal/repository"
)

type Authorization interface {
	Test()
	CreateUser(user core.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

// Главный тип слоя SVC, который встраивается как зависимость в слое TRANSPORT
type Service struct {
	Authorization
}

// Конструктор слоя SVC
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
