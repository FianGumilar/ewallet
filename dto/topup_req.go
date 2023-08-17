package dto

type TopUpRequest struct {
	Amount float64 `json:"amount"`
	UserID int64   `json:"user_id"`
}
