package project

type StartRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Organisation string `json:"organisation"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
}
