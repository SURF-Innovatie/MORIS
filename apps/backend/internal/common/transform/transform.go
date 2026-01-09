package transform

import "github.com/samber/lo"

// ToEntity transforms a single Ent row to a Domain Entity.
// It handles the pointer receiver requirement for FromEnt using a generic type constraint P.
// P must be a pointer to E, and must implement FromEnt.
// Usage: ToEntity[entities.User](userRow)
func ToEntity[E any, P interface {
	*E
	FromEnt(*EntRow) *E
}, EntRow any](row *EntRow) E {
	var e E
	return *P(&e).FromEnt(row)
}

// ToEntities transforms a slice of Ent rows to a slice of Domain Entities.
// Usage: ToEntities[entities.User](userRows)
func ToEntities[E any, P interface {
	*E
	FromEnt(*EntRow) *E
}, EntRow any](rows []*EntRow) []E {
	return lo.Map(rows, func(row *EntRow, _ int) E {
		var e E
		return *P(&e).FromEnt(row)
	})
}

// ToEntitiesPtr transforms a slice of Ent rows to a slice of pointers to Domain Entities.
// Usage: ToEntitiesPtr[entities.User](userRows)
func ToEntitiesPtr[E any, P interface {
	*E
	FromEnt(*EntRow) *E
}, EntRow any](rows []*EntRow) []*E {
	return lo.Map(rows, func(row *EntRow, _ int) *E {
		var e E
		return P(&e).FromEnt(row)
	})
}

// ToEntityPtr transforms a single Ent row to a pointer to a Domain Entity.
// Usage: ToEntityPtr[entities.User](userRow)
func ToEntityPtr[E any, P interface {
	*E
	FromEnt(*EntRow) *E
}, EntRow any](row *EntRow) *E {
	var e E
	return P(&e).FromEnt(row)
}

// ToDTOItem transforms a single Domain Entity to a DTO.
// Usage: ToDTOItem[dto.UserResponse](userEntity)
func ToDTOItem[DTO interface{ FromEntity(E) DTO }, E any](entity E) DTO {
	var d DTO
	return d.FromEntity(entity)
}

// ToDTOs transforms a slice of Domain Entities to a slice of DTOs.
// Usage: ToDTOs[dto.UserResponse](userEntities)
func ToDTOs[DTO interface{ FromEntity(E) DTO }, E any](entities []E) []DTO {
	return lo.Map(entities, func(entity E, _ int) DTO {
		var d DTO
		return d.FromEntity(entity)
	})
}
