package model

type Orders struct {
	OrderID    int     `json:"order_id"`
	UserID     int     `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
}
