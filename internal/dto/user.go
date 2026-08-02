package dto

import "time"

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email" opt:"true"`
	Username string `json:"username" opt:"true"`
	Password string `json:"password"`
}

type EditUserAccountRequest struct {
	NewUsername   string `json:"username" opt:"true"`
	NewPassword   string `json:"password" opt:"true"`
	DeleteAccount bool   `json:"delete" opt:"true"`
}

type GetTopUsersRequest struct {
	Cursor int `json:"cursor"`
	Limit  int `json:"limit"`
}
type GetTopUsersResponse struct {
	Usernames  []string `json:"usernames"`
	TotalTimes []int    `json:"total_times"`
	TodayTimes []int    `json:"today_times"`
}

type GetUserAccountResponse struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
