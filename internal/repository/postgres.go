package repository

import (
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	usersTable = "users"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(urlDB string) (*sqlx.DB, error) {
	if urlDB == "" {
		return nil, errors.New("url not set")
	}
	db, err := sqlx.Open("postgres", urlDB)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	logrus.Info("Success connect to postgres.")

	// делаем запрос
	var checkExist bool
	row := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE  tablename  = 'users');")
	err = row.Scan(&checkExist)
	if err != nil {
		log.Fatal(err)
	}
	if !checkExist {
		_, err = db.Exec("CREATE TABLE users (id serial not null unique,login varchar(255) not null unique, password_hash varchar(255) not null);") //QueryRowContext т.к. одна запись
		if err != nil {
			log.Fatal(err)
		}
		logrus.Info("Table users successful create")

	} else {
		logrus.Info("Table users already created")
	}

	return db, nil
}

// "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"
