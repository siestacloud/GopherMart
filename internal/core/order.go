package core

type Order struct {
	ID         string `json:"number,omitempty" db:"user_order" validate:"required,numeric,max=20"`
	Status     string `json:"status" db:"status"`
	CreateTime string `json:"uploaded_at" db:"create_time"`
}
