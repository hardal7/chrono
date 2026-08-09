package dto

type CreateTopicRequest struct {
	Name string `json:"name"`
}

type EditTopicRequest struct {
	Name    string `json:"name"`
	NewName string `json:"new_name" validate:"omitempty"`
	Delete  bool   `json:"delete" validate:"omitempty"`
}

type GetTopicRequest struct {
	Name string `json:"name"`
}
type GetTopicResponse struct {
	TotalTime int `json:"total_time_tracked_seconds"`
}
