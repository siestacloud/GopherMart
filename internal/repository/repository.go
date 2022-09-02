package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/gopherMart/internal/core"
)

type Authorization interface {
	TestDB()
	CreateUser(user core.User) (int, error)
	GetUser(username, password string) (core.User, error)
}

// Главный тип слоя repository, который встраивается как зависимость в слое SVC
type Repository struct {
	Authorization
}

//Конструктор слоя repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
