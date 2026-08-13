package dto

type CreateTopicRequest struct {
	Name string `json:"name"`
}

type EditTopicRequest struct {
	Name    string `json:"name"`
	NewName string `json:"new_name" validate:"omitempty"`
	Delete  bool   `json:"delete" validate:"omitempty"`
}

type GetTopicNamedRequest struct {
	Name string `json:"name"`
}
type GetTopicNamedResponse struct {
	TotalTime int `json:"total_time_tracked_seconds"`
	Streak    int `json:"streak"`
}

type TopicSelection struct {
	Name      string `json:"name"`
	TotalTime int    `json:"total_time_tracked_seconds"`
}
type GetTopicsAllResponse struct {
	Topics []TopicSelection `json:"topics"`
}
