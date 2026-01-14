package dto

import (
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type UserPaginatedResponse struct {
	Data       []UserResponse `json:"data"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

type UserRequest struct {
	PersonID   uuid.UUID `json:"person_id"`
	Password   string    `json:"password"`
	IsSysAdmin *bool     `json:"is_sys_admin"`
}

type UserToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	PersonID    uuid.UUID `json:"person_id"`
	ORCiD       *string   `json:"orcid"`
	Name        string    `json:"name"`
	GivenName   *string   `json:"givenName"`
	FamilyName  *string   `json:"familyName"`
	Email       string    `json:"email"`
	AvatarURL   *string   `json:"avatarUrl"`
	Description *string   `json:"description"`
	IsSysAdmin  bool      `json:"is_sys_admin"`
	IsActive    bool      `json:"is_active"`
}

func (r UserResponse) FromEntity(acc *entities.UserAccount) UserResponse {
	p := acc.Person
	u := acc.User

	return UserResponse{
		ID:          u.ID,
		PersonID:    u.PersonID,
		Email:       p.Email,
		Name:        p.Name,
		ORCiD:       p.ORCiD,
		GivenName:   p.GivenName,
		FamilyName:  p.FamilyName,
		AvatarURL:   p.AvatarUrl,
		Description: p.Description,
		IsSysAdmin:  u.IsSysAdmin,
		IsActive:    u.IsActive,
	}
}

type UserPersonResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	GivenName  *string   `json:"givenName"`
	FamilyName *string   `json:"familyName"`
	Email      string    `json:"email"`
	AvatarURL  *string   `json:"avatarUrl"`
	ORCiD      *string   `json:"orcid"`
}

func (r UserPersonResponse) FromEntity(p entities.Person) UserPersonResponse {
	return UserPersonResponse{
		ID:         p.ID,
		Name:       p.Name,
		GivenName:  p.GivenName,
		FamilyName: p.FamilyName,
		Email:      p.Email,
		AvatarURL:  p.AvatarUrl,
		ORCiD:      p.ORCiD,
	}
}
