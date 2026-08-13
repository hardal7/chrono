package dto

type CreateFriendRequestRequest struct {
	Username string `json:"username"`
}

type AcceptFriendRequestRequest struct {
	Username string `json:"username"`
}

type RemoveFriendRequest struct {
	Username string `json:"username"`
}
