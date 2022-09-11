package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// BalancePostgres
type BalancePostgres struct {
	db *sqlx.DB
}

// NewBalancePostgres
func NewBalancePostgres(db *sqlx.DB) *BalancePostgres {
	return &BalancePostgres{
		db: db,
	}
}

// Create транзакция. Создаю баланс в базу и связывую с новым клиентом
// * метод используется при авторизации нового клиента
func (o *BalancePostgres) Create(userId int) error {
	if o.db == nil {
		return errors.New("database are not connected")
	}
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	var id int
	balanceQuery := fmt.Sprintf("INSERT INTO %s (current, withdrawn) VALUES ($1,$2) RETURNING id", balanceTable)
	row := tx.QueryRow(balanceQuery, 0, 0)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	usersBalanceQuery := fmt.Sprintf("INSERT INTO %s (user_id, balance_id) VALUES ($1, $2)", userBalanceTable)
	_, err = tx.Exec(usersBalanceQuery, userId, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}
