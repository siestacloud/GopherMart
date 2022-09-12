package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/gopherMart/pkg"
	"github.com/sirupsen/logrus"
)

const (
	usersTable        = "users"
	ordersTable       = "orders"
	balanceTable      = "balance"
	withdrawTable     = "withdraw"
	userOrderTable    = "users_orders"
	userBalanceTable  = "users_balance"
	userWithdrawTable = "users_withdraw"
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
	if err := createTable(db, usersTable, "CREATE TABLE users (id serial not null unique,login varchar(255) not null unique, password_hash varchar(255) not null);"); err != nil {
		log.Fatal(err)
	}
	// делаем запрос
	if err := createTable(db, ordersTable, "CREATE TABLE orders (id serial not null unique,user_order bigint not null unique, status varchar(255) not null,sum numeric not null,create_time timestamp,withdrawn_time timestamp NULL);"); err != nil {
		log.Fatal(err)
	}
	// делаем запрос
	if err := createTable(db, userOrderTable, "CREATE TABLE users_orders (id serial not null unique,user_id int references users (id) on delete cascade not null,order_id int references orders (id) on delete cascade not null);"); err != nil {
		log.Fatal(err)
	}
	// делаем запрос
	if err := createTable(db, balanceTable, "CREATE TABLE balance (id serial not null unique,current numeric, withdrawn numeric);"); err != nil {
		log.Fatal(err)
	}
	// делаем запрос
	if err := createTable(db, userBalanceTable, "CREATE TABLE users_balance (id serial not null unique,user_id int references users (id) on delete cascade not null,balance_id int references balance (id) on delete cascade not null);"); err != nil {
		log.Fatal(err)
	}
	return db, nil
}

// * "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"

func createTable(db *sqlx.DB, nameTable, query string) error {

	var checkExist bool

	row := db.QueryRow(fmt.Sprintf("SELECT EXISTS (SELECT FROM pg_tables WHERE  tablename  = '%s');", nameTable))
	if err := row.Scan(&checkExist); err != nil {
		return err
	}

	if !checkExist {
		_, err := db.Exec(query) //QueryRowContext т.к. одна запись
		if err != nil {
			return err
		}
		pkg.InfoPrint("repository", "ok", "Table  successful create")

	} else {
		pkg.WarnPrint("repository", "ok", "Table  already created")
	}

	return nil
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
