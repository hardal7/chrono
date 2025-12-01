package model

import "time"

type Session struct {
	ID        int           `db:"id"`
	Name      string        `db:"name"`
	Password  string        `db:"password"`
	Admin     User          `db:"admin"`
	Users     []SessionUser `db:"users"`
	Expiry    time.Time     `db:"expiry"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}

type SessionUser struct {
	ID     int `db:"id"`
	UserID int `db:"userid"`
}

type CreateSessionRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Expiry   int    `json:"expiry"`
}

type EditSessionRequest struct {
	NewName       string        `json:"name"`
	NewPassword   string        `json:"password"`
	NewExpiry     int           `json:"expiry"`
	RemovedUsers  []SessionUser `json:"removed_users"`
	DeleteSession int           `json:"delete"`
}
