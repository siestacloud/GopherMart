package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
	TestDB()
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
