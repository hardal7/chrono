package dto

import "time"

type TrackTopicEventRequest struct {
	Topic       string    `json:"topic"`
	TimeSeconds int       `json:"time_seconds"`
	Date        time.Time `json:"date"`
}

type GetTopicEventsRequest struct {
	Topic    string    `json:"topic" validate:"omitempty"`
	FromDate time.Time `json:"from_date" validate:"omitempty"`
	ToDate   time.Time `json:"to_date" validate:"omitempty"`
}
type GetTopicEventsResponse struct {
	Topics       []string    `json:"topics"`
	TimesTracked []int       `json:"times_tracked"`
	Dates        []time.Time `json:"dates"`
}

type GetTopicEventsTodayResponse struct {
	TotalTime int `json:"total_time"`
}
