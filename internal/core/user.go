package core

//User имплементирует клиента
type User struct {
	ID       int    `json:"-" db:"id"`
	Login    string `json:"login"  validate:"required"`
	Password string `json:"password" validate:"required"`
}
