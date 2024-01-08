package model

type Users struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	JWTToken string `json:"jwt_token"`
}
