package core

type Order struct {
	ID int `db:"id" validate:"required"`
}
