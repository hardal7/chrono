package topicevent

import "time"

type TopicEvent struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	TopicID     int       `db:"topic_id"`
	TimeTracked int       `db:"time_tracked_seconds"`
	Date        time.Time `db:"date"`
}

type TrackRequest struct {
	Topic       string    `json:"topic"`
	TimeSeconds int       `json:"time"`
	Date        time.Time `json:"date"`
}

type GetRequest struct {
	Topic string `json:"topic" opt:"true"`
	Date  []int  `json:"date" opt:"true"`
}

type GetResponse struct {
	Topics       []string    `json:"topics"`
	TimesTracked []int       `json:"times_tracked"`
	Dates        []time.Time `json:"dates"`
}

type GetTodayResponse struct {
	TotalTime int `json:"total_time"`
}
