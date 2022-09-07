package core

type Order struct {
	ID         int    `json:"id,omitempty" db:"id" validate:"required"`
	UserOrder  int64  `json:"number" db:"user_order"`
	Status     string `json:"status" db:"status"`
	CreateTime string `json:"uploaded_at" db:"create_time"`
}
