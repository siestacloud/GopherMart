package core

// Order имплементирует заказ клиента
type Order struct {
	Number     string `json:"number,omitempty" db:"user_order" validate:"required,numeric,max=20"`
	Status     string `json:"status" db:"status"`
	Accrual    int32  `json:"accrual"`
	CreateTime string `json:"uploaded_at" db:"create_time"`
}
