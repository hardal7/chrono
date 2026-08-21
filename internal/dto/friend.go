package dto

import "time"

type CreateFriendRequestRequest struct {
	Username string `json:"username"`
}

type AcceptFriendRequestRequest struct {
	Username string `json:"username"`
}

type RemoveFriendRequest struct {
	Username string `json:"username"`
}

type FriendRequest struct {
	FromUsername string    `json:"from_username"`
	Date         time.Time `json:"date"`
}
type GetFriendRequestsAllResponse struct {
	Requests []FriendRequest `json:"requests"`
}
