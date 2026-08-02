package dto

import "time"

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"email,omitempty"`
	Username string `json:"username" validate:"omitempty"`
	Password string `json:"password"`
}

type EditUserAccountRequest struct {
	NewUsername   string `json:"username" validate:"omitempty"`
	NewPassword   string `json:"password" validate:"omitempty"`
	DeleteAccount bool   `json:"delete" validate:"omitempty"`
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
	Email     string    `json:"email" validate:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
