package dto

type CreateTopicRequest struct {
	Name string `json:"name" validate:"max=64"`
}

type EditTopicRequest struct {
	Name    string `json:"name"`
	NewName string `json:"new_name" validate:"max=64"`
}

type DeleteTopicRequest struct {
	Name string `json:"name"`
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
