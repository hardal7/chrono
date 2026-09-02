package dto

import "time"

type CreateSessionRequest struct {
	Name            string    `json:"name" validate:"max=64"`
	MaxParticipants int       `json:"max_participants" validate:"omitempty,max=1024"`
	ExpiresAt       time.Time `json:"expires_at" validate:"omitempty"`
	Topic           string    `json:"topic" validate:"omitempty"`
}

type EditSessionRequest struct {
	Name               string    `json:"name"`
	NewName            string    `json:"new_name" validate:"omitempty"`
	NewMaxParticipants int       `json:"max_participants" validate:"omitempty"`
	NewExpiresAt       time.Time `json:"expires_at" validate:"omitempty"`
}

type DeleteSessionRequest struct {
	Name string `json:"name"`
}

type JoinSessionRequest struct {
	Name          string `json:"name"`
	OwnerUsername string `json:"owner_username"`
}

type KickFromSessionRequest struct {
	SessionName         string `json:"session_name"`
	ParticipantUsername string `json:"participant_username"`
}

type MinParticipant struct {
	Name       string `json:"name"`
	AvatarPath string `json:"avatar_path"`
}
type SessionSelection struct {
	Name              string           `json:"name"`
	OwnerUsername     string           `json:"owner_username"`
	OwnerAvatarPath   string           `json:"owner_avatar_path"`
	ExpiresAt         *time.Time       `json:"expires_at"`
	TotalTime         int              `json:"total_time_seconds"`
	TotalParticipants int              `json:"total_participants"`
	MaxParticipants   int              `json:"max_participants"`
	Participants      []MinParticipant `json:"participants"`
}
type GetSessionsAllResponse struct {
	Sessions []SessionSelection `json:"sessions"`
}

type GetSessionNamedRequest struct {
	Name          string `json:"name"`
	OwnerUsername string `json:"owner_username"`
}
type Participant struct {
	Name             string `json:"name"`
	AvatarPath       string `json:"avatar_path"`
	SessionTime      int    `json:"session_time_tracked_seconds"`
	SessionTimeToday int    `json:"session_time_tracked_today_seconds"`
	LastOnline       int    `json:"last_online_seconds_ago"`
}
type GetSessionNamedResponse struct {
	Name              string        `json:"name"`
	OwnerUsername     string        `json:"owner_username"`
	ExpiresAt         *time.Time    `json:"expires_at"`
	TotalTime         int           `json:"total_time_tracked_seconds"`
	TotalParticipants int           `json:"total_participants"`
	MaxParticipants   int           `json:"max_participants"`
	Participants      []Participant `json:"participants"`
}
