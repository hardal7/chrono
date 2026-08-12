package dto

import "time"

type CreateSessionRequest struct {
	Name            string    `json:"name"`
	MaxParticipants int       `json:"max_participants" validate:"omitempty"`
	Password        string    `json:"password" validate:"omitempty"`
	ExpiresAt       time.Time `json:"expires_at" validate:"omitempty"`
	Topic           string    `json:"topic" validate:"omitempty"`
}

type EditSessionRequest struct {
	Name               string    `json:"name"`
	NewName            string    `json:"new_name" validate:"omitempty"`
	NewMaxParticipants int       `json:"max_participants" validate:"omitempty"`
	NewPassword        string    `json:"password" validate:"omitempty"`
	NewExpiresAt       time.Time `json:"expires_at" validate:"omitempty"`
	KickedUsername     string    `json:"kicked_username" validate:"omitempty"`
	Delete             bool      `json:"delete" validate:"omitempty"`
}

type JoinSessionRequest struct {
	Name          string `json:"name"`
	OwnerUsername string `json:"owner_username"`
	Password      string `json:"password" validate:"omitempty"`
}
