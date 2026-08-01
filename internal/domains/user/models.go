package user

import "time"

type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	TotalTime int       `db:"total_time_tracked_seconds"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" opt:"true"`
	Username string `json:"username" opt:"true"`
	Password string `json:"password"`
}

type EditAccountRequest struct {
	NewUsername   string `json:"username" opt:"true"`
	NewPassword   string `json:"password" opt:"true"`
	DeleteAccount bool   `json:"delete" opt:"true"`
}

type GetAccountResponse struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type GetTopUsersRequest struct {
	Cursor int `json:"cursor"`
	Limit  int `json:"limit"`
}
