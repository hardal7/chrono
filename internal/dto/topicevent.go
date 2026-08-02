package dto

import "time"

type TrackTopicEventRequest struct {
	Topic       string    `json:"topic"`
	TimeSeconds int       `json:"time_seconds"`
	Date        time.Time `json:"date"`
}

type GetTopicEventsRequest struct {
	Topic string      `json:"topic" opt:"true"`
	Date  []time.Time `json:"date" opt:"true"`
}
type GetTopicEventsResponse struct {
	Topics       []string    `json:"topics"`
	TimesTracked []int       `json:"times_tracked"`
	Dates        []time.Time `json:"dates"`
}

type GetTopicEventsTodayResponse struct {
	TotalTime int `json:"total_time"`
}
