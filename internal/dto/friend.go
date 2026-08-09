package dto

type CreateFriendRequestRequest struct {
	Username string `json:"username"`
}

type AcceptFriendRequestRequest struct {
	Username string `json:"username"`
}

type DeleteFriendRequest struct {
	Username string `json:"username"`
}
