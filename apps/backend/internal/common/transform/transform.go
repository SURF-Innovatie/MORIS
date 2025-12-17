package transform

// Map transforms a slice of Source items to a slice of Target items using a mapping function.
// This is useful for general transformations.
func Map[S any, T any](source []S, mapper func(S) T) []T {
	result := make([]T, len(source))
	for i, item := range source {
		result[i] = mapper(item)
	}
	return result
}

// ToEntity transforms a single Ent row to a Domain Entity.
// It handles the pointer receiver requirement for FromEnt using a generic type constraint P.
// P must be a pointer to E, and must implement FromEnt.
// Usage: ToEntity[entities.User](userRow)
func ToEntity[E any, P interface {
	*E
	FromEnt(*EntRow) *E
}, EntRow any](row *EntRow) E {
	var e E
	// P(&e) converts *E to P, so we can call FromEnt on it.
	// However, FromEnt usually returns *E (newly allocated or same).
	// If FromEnt is defined on *E, we can use P.FromEnt.
	// We need to create an instance of P to call FromEnt if it is a method on the type.
	// But FromEnt is usually a factory-like method or a converter on a zero value.
	// Looking at existing usage (u.FromEnt(row)), it seems to be used on an instance.
	// Let's create a zero value of E, take its address (which satisfies P), and call FromEnt.
	return *P(&e).FromEnt(row)
}

// ToEntities transforms a slice of Ent rows to a slice of Domain Entities.
// Usage: ToEntities[entities.User](userRows)
func ToEntities[E any, P interface {
	*E
	FromEnt(*EntRow) *E
}, EntRow any](rows []*EntRow) []E {
	result := make([]E, len(rows))
	for i, row := range rows {
		var e E
		result[i] = *P(&e).FromEnt(row)
	}
	return result
}

// ToEntitiesPtr transforms a slice of Ent rows to a slice of pointers to Domain Entities.
// Usage: ToEntitiesPtr[entities.User](userRows)
func ToEntitiesPtr[E any, P interface {
	*E
	FromEnt(*EntRow) *E
}, EntRow any](rows []*EntRow) []*E {
	result := make([]*E, len(rows))
	for i, row := range rows {
		var e E
		result[i] = P(&e).FromEnt(row)
	}
	return result
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
	result := make([]DTO, len(entities))
	for i, entity := range entities {
		var d DTO
		result[i] = d.FromEntity(entity)
	}
	return result
}
