package core

// Order имплементирует баланс клиента
type Balance struct {
	ID        int     `db:"id"`                                   // * уникальный идентификатор баланса клиента
	Current   float64 `json:"current,omitempty" db:"current"`     // * Текущее количество баллов клиента
	Withdrawn float64 `json:"withdrawn,omitempty" db:"withdrawn"` // * Общее количество использованных баллов за все время
}
