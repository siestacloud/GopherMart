package core

// Order имплементирует заказ клиента
type Order struct {
	Number     string  `json:"number,omitempty" db:"user_order" validate:"required,numeric,max=20"`
	Status     string  `json:"status,omitempty" db:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	CreateTime string  `json:"uploaded_at,omitempty" db:"create_time"`
}
