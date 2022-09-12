package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/siestacloud/gopherMart/internal/core"
	"github.com/siestacloud/gopherMart/pkg"
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

// Create транзакция. Создаю баланс в базе и связывую с новым клиентом
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

// Get получить текущее количество баллов клиента и общее количество использованных баллов за все время
func (m *BalancePostgres) Get(userID int) (*core.Balance, error) {
	if m.db == nil {
		return nil, errors.New("database are not connected")
	}

	var balance core.Balance
	query := fmt.Sprintf(`SELECT balance_id FROM %s  WHERE user_id = $1`, userBalanceTable)
	if err := m.db.Get(&balance.ID, query, userID); err != nil {
		pkg.ErrPrint("repository", 500, err)
		return nil, err
	}

	query = fmt.Sprintf(`SELECT current,withdrawn FROM %s  WHERE id = $1`, balanceTable)
	if err := m.db.Get(&balance, query, balance.ID); err != nil {
		pkg.ErrPrint("repository", 500, err)
		return nil, err
	}
	return &balance, nil
}

// UpdateCurrent обновить текущее количество баллов клиента
func (m *BalancePostgres) UpdateCurrent(balance *core.Balance) error {
	if m.db == nil {
		return errors.New("database are not connected")
	}
	balanceQuery := fmt.Sprintf("UPDATE %s SET current = %v WHERE id = %v ", balanceTable, balance.Current, balance.ID)
	_, err := m.db.Exec(balanceQuery)
	if err != nil {
		pkg.ErrPrint("repository", 500, err)
		return err
	}

	return nil
}

// UpdateWithdrawn обновить общее количество использованных баллов клиента за все время
func (m *BalancePostgres) UpdateWithdrawn(balance *core.Balance) error {
	if m.db == nil {
		return errors.New("database are not connected")
	}
	balanceQuery := fmt.Sprintf("UPDATE %s SET withdrawn = %v WHERE id = %v ", balanceTable, balance.Withdrawn, balance.ID)
	_, err := m.db.Exec(balanceQuery)
	if err != nil {
		pkg.ErrPrint("repository", 500, err)
		return err
	}

	return nil
}
