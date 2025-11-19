package person

type CreatePersonRequest struct {
	Name       string  `json:"name"`
	GivenName  *string `json:"givenName"`
	FamilyName *string `json:"familyName"`
	Email      *string `json:"email"`
}

type PersonResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	GivenName  *string `json:"givenName"`
	FamilyName *string `json:"familyName"`
	Email      *string `json:"email"`
}
