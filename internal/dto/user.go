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
type TopUser struct {
	Rank      int    `json:"rank"`
	Username  string `json:"username"`
	TotalTime int    `json:"total_time"`
	TodayTime int    `json:"today_time"`
}
type GetTopUsersResponse struct {
	Users []TopUser `json:"users"`
}

type GetUserAccountResponse struct {
	Email     string    `json:"email" validate:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
