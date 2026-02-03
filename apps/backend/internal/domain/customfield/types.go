package customfield

type Type string

const (
	TypeText    Type = "TEXT"
	TypeNumber  Type = "NUMBER"
	TypeBoolean Type = "BOOLEAN"
	TypeDate    Type = "DATE"
)

func (t Type) Valid() bool {
	switch t {
	case TypeText, TypeNumber, TypeBoolean, TypeDate:
		return true
	default:
		return false
	}
}
