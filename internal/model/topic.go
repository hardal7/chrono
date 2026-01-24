package model

import "time"

type Topic struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedBy int       `db:"created_by_userid"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type TopicEvent struct {
	ID          int `db:"id"`
	UserID      int `db:"user_id"`
	TopicID     int `db:"topic_id"`
	TimeTracked int `db:"time_tracked"`
	Date        int `db:"date"`
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

type TrackTopicRequest struct {
	Topic string    `json:"topic"`
	Time  time.Time `json:"time"`
	Date  time.Time `json:"date"`
}

type TopicEventRequest struct {
	Topic string `json:"topic"`
	Date  []int  `json:"date"`
}

type TopicEventResponse struct {
	Topics       []string `json:"topics"`
	TimesTracked []int    `json:"times_tracked"`
	Dates        []int    `json:"dates"`
}
