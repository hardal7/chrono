package model

import "time"

type Topic struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedBy int       `db:"created_by_userid"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type TopicUser struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	TopicID     int       `db:"topic_id"`
	TimeTracked time.Time `db:"time_tracked"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CreateTopicRequest struct {
	ID   int    `db:"id"`
	Name string `json:"name"`
}

type EditTopicRequest struct {
	Name        string `json:"name"`
	NewName     string `json:"new_name"`
	DeleteTopic bool   `json:"delete"`
}
