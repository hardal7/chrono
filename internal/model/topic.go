package model

import "time"

type Topic struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type TopicUser struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	TagID       int       `db:"topic_id"`
	TimeTracked time.Time `db:"time_tracked"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateTopicRequest struct {
	Name string `json:"name"`
}

type EditTopicRequest struct {
	NewName   string `json:"name"`
	DeleteTag bool   `json:"delete"`
}
