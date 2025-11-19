package project

// StartRequest represents the request body for starting a new project
type StartRequest struct {
	Title        string `json:"title" example:"New Project"`
	Description  string `json:"description" example:"This is a new project"`
	Organisation string `json:"organisation" example:"SURF"`
	StartDate    string `json:"startDate" example:"2025-01-01T00:00:00Z"`
	EndDate      string `json:"endDate" example:"2025-12-31T23:59:59Z"`
}
