package odrl

// Operator represents the comparison operator.
type Operator string

// Common ODRL Operators
const (
	OperatorEq       Operator = "eq"
	OperatorGt       Operator = "gt"
	OperatorGteq     Operator = "gteq"
	OperatorLt       Operator = "lt"
	OperatorLteq     Operator = "lteq"
	OperatorNeq      Operator = "neq"
	OperatorIsAny    Operator = "isAny"
	OperatorIsAll    Operator = "isAll"
	OperatorIsNone   Operator = "isNone"
	OperatorIsPartOf Operator = "isPartOf"
)

// Constraint represents a boolean expression.
type Constraint struct {
	LeftOperand  string      `json:"leftOperand"`
	Operator     Operator    `json:"operator"`
	RightOperand interface{} `json:"rightOperand,omitempty"` // Can be literal, IRI, or object
	Unit         string      `json:"unit,omitempty"`
	DataType     string      `json:"dataType,omitempty"`
}

// NewConstraint creates a new constraint.
func NewConstraint(leftOperand string, operator Operator, rightOperand interface{}) *Constraint {
	return &Constraint{
		LeftOperand:  leftOperand,
		Operator:     operator,
		RightOperand: rightOperand,
	}
}

// WithUnit sets the unit of the constraint.
func (c *Constraint) WithUnit(unit string) *Constraint {
	c.Unit = unit
	return c
}

// WithDataType sets the data type of the constraint.
func (c *Constraint) WithDataType(dataType string) *Constraint {
	c.DataType = dataType
	return c
}
