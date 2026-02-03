package customfield

type Category string

const (
	CategoryProject Category = "PROJECT"
	CategoryPerson  Category = "PERSON"
)

func (c Category) Valid() bool {
	switch c {
	case CategoryProject, CategoryPerson:
		return true
	default:
		return false
	}
}
