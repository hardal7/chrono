package dto

import "time"

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"email"`
	Username string `json:"username" validate:"min=4,max=32"`
	Password string `json:"password" validate:"min=4"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty"`
	Password string `json:"password"`
}

type EditUserAccountRequest struct {
	NewUsername string `json:"username" validate:"omitempty,max=32"`
	NewPassword string `json:"password" validate:"omitempty"`
}

type ResetUserPasswordRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=4,max=32"`
}

type GetTopUsersRequest struct {
	Cursor    int    `json:"cursor"`
	Limit     int    `json:"limit"`
	Scope     string `json:"scope" validate:"omitempty,oneof=friends local global"`
	MatchName string `json:"match_name"`
}
type TopUser struct {
	Rank       int    `json:"rank"`
	Username   string `json:"username"`
	TotalTime  int    `json:"total_time"`
	TodayTime  int    `json:"today_time"`
	AvatarPath string `json:"avatar_path"`
}
type GetTopUsersResponse struct {
	Users []TopUser `json:"users"`
}

type GetUserAccountResponse struct {
	Email     string    `json:"email" validate:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type GetUserProfileResponse struct {
	Username         string `json:"username"`
	AvatarPath       string `json:"avatar_path"`
	Country          string `json:"country"`
	TotalTimeSeconds int    `json:"total_time_seconds"`
	TodayTimeSeconds int    `json:"today_time_seconds"`
	BestTopic        string `json:"best_topic"`
	FriendStatus     string `json:"friend_status"`
}
