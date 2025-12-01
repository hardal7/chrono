package model

import "time"

type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EditAccountRequest struct {
	NewUsername   string `json:"username"`
	NewPassword   string `json:"password"`
	DeleteAccount bool   `json:"delete_account"`
}
