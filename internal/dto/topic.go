package dto

type CreateTopicRequest struct {
	ID   int    `db:"id"`
	Name string `json:"name"`
}

type EditTopicRequest struct {
	Name    string `json:"name"`
	NewName string `json:"new_name" opt:"true"`
	Delete  bool   `json:"delete" opt:"true"`
}

type GetTopicRequest struct {
	Name string `json:"name"`
}
type GetTopicResponse struct {
	TotalTime int `json:"total_time_tracked_seconds"`
}
