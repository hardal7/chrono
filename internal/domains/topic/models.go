package topic

import "time"

type Topic struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	TotalTime int       `db:"total_time_tracked_seconds"`
	CreatedBy int       `db:"created_by_userid"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateRequest struct {
	ID   int    `db:"id"`
	Name string `json:"name"`
}

type EditRequest struct {
	Name        string `json:"name"`
	NewName     string `json:"new_name" opt:"true"`
	DeleteTopic bool   `json:"delete" opt:"true"`
}

type GetRequest struct {
	Name string `json:"name"`
}

type GetResponse struct {
	TotalTime int `json:"total_time_tracked_seconds"`
}
