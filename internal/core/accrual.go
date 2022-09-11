package core

import "time"

// AccrualAPIError для парсинга ответа (err) от накопительной системы лояльности
type AccrualAPIError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
