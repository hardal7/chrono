package dto

type CreateSessionRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Expiry   int    `json:"expiry"`
}

type JoinSessionRequest struct {
	SessionName   string `json:"session_name"`
	OwnerUsername string `json:"owner_username"`
	Password      string `json:"password"`
}

type EditSessionRequest struct {
	Name          string   `json:"name"`
	NewName       string   `json:"new_name"`
	NewPassword   string   `json:"password"`
	NewExpiry     int      `json:"expiry"`
	RemovedUsers  []string `json:"removed_users"`
	DeleteSession int      `json:"delete"`
}
