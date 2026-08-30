package dto

type CreateFeatureRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Title    string `json:"title" validate:"required"`
	Problem  string `json:"problem" validate:"required"`
	Feature  string `json:"feature" validate:"required"`
	Priority string `json:"priority" validate:"omitempty,oneof=low medium high"`
}

type CreateBugReport struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Steps       string `json:"steps"`
	Environment string `json:"environment"`
	Additional  string `json:"additional"`
}
