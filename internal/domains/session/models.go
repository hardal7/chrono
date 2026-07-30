package session

import "time"

type Session struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	Admin     int       `db:"admin_id"`
	Users     []int     `db:"user_ids"`
	Expiry    time.Time `db:"expiry"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Expiry   int    `json:"expiry"`
}

type JoinRequest struct {
	SessionID string `json:"session_id"`
	Password  string `json:"password"`
}

type EditRequest struct {
	SessionID     int      `json:"id"`
	NewName       string   `json:"name"`
	NewPassword   string   `json:"password"`
	NewExpiry     int      `json:"expiry"`
	RemovedUsers  []string `json:"removed_users"`
	DeleteSession int      `json:"delete"`
}
