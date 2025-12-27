package events

type Status string

const (
	StatusApproved Status = "approved"
	StatusPending  Status = "pending"
	StatusRejected Status = "rejected"
)
