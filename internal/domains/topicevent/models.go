package topicevent

import "time"

type TopicEvent struct {
	ID          int `db:"id"`
	UserID      int `db:"user_id"`
	TopicID     int `db:"topic_id"`
	TimeTracked int `db:"time_tracked"`
	Date        int `db:"date"`
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
