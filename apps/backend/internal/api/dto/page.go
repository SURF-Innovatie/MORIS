package dto

import (
	"time"

	"github.com/SURF-Innovatie/MORIS/internal/types"
	"github.com/google/uuid"
)

type CreatePageRequest struct {
	Title       string          `json:"title"`
	Slug        string          `json:"slug"`
	Type        string          `json:"type" example:"project"` // project, profile
	Content     []types.Section `json:"content"`
	ProjectID   *uuid.UUID      `json:"project_id"`
	UserID      *uuid.UUID      `json:"user_id"`
	IsPublished bool            `json:"is_published"`
}

type UpdatePageRequest struct {
	Title       *string         `json:"title"`
	Content     []types.Section `json:"content"`
	IsPublished *bool           `json:"is_published"`
}

type PageResponse struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Slug        string          `json:"slug"`
	Type        string          `json:"type"`
	Content     []types.Section `json:"content"`
	IsPublished bool            `json:"is_published"`
	ProjectID   *uuid.UUID      `json:"project_id,omitempty"`
	UserID      *uuid.UUID      `json:"user_id,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
