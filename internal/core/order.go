package core

// Order имплементирует заказ клиента (в счет накопления и в счет списания баллов)
type Order struct {
	Number        string  `json:"number,omitempty" db:"user_order" validate:"required,numeric,max=20"`
	Status        string  `json:"status,omitempty" db:"status"`
	Accrual       float64 `json:"accrual,omitempty"`
	Sum           float64 `json:"sum,omitempty" db:"sum"`
	CreateTime    string  `json:"uploaded_at,omitempty" db:"update_time"`
	WithdrawnTime string  `json:"processed_at,omitempty" `
}
