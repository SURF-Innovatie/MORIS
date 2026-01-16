package entities

// ODataQuery represents a parsed OData request
type ODataQuery struct {
	Select  []string
	Filter  *ODataFilter
	OrderBy []ODataOrderBy
	Top     *int
	Skip    *int
	Expand  []string
	Count   bool // Request count in response
}

// ODataFilter represents filter expressions
type ODataFilter struct {
	Field    string
	Operator ODataOperator
	Value    any
	And      *ODataFilter
	Or       *ODataFilter
}

// ODataOperator represents comparison operators
type ODataOperator string

const (
	ODataOpEqual          ODataOperator = "eq"
	ODataOpNotEqual       ODataOperator = "ne"
	ODataOpGreaterThan    ODataOperator = "gt"
	ODataOpLessThan       ODataOperator = "lt"
	ODataOpGreaterOrEqual ODataOperator = "ge"
	ODataOpLessOrEqual    ODataOperator = "le"
	ODataOpContains       ODataOperator = "contains"
	ODataOpStartsWith     ODataOperator = "startswith"
	ODataOpEndsWith       ODataOperator = "endswith"
)

// ODataOrderBy represents ordering
type ODataOrderBy struct {
	Field string
	Desc  bool
}

// ODataResult wraps paginated results with OData metadata
type ODataResult[T any] struct {
	Value    []T     `json:"value"`
	Count    *int    `json:"@odata.count,omitempty"`
	NextLink *string `json:"@odata.nextLink,omitempty"`
	Context  string  `json:"@odata.context,omitempty"`
}

// NewODataResult creates a new OData result wrapper
func NewODataResult[T any](value []T, count *int, nextLink *string) ODataResult[T] {
	return ODataResult[T]{
		Value:    value,
		Count:    count,
		NextLink: nextLink,
	}
}

// ODataError represents an OData-formatted error response
type ODataError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ODataErrorResponse wraps an OData error
type ODataErrorResponse struct {
	Error ODataError `json:"error"`
}

// AllowedFields defines which fields can be used in OData queries for security
type AllowedFields struct {
	Selectable []string
	Filterable []string
	Sortable   []string
	Expandable []string
}

// BudgetODataFields defines allowed OData fields for budget queries
var BudgetODataFields = AllowedFields{
	Selectable: []string{"id", "projectId", "title", "status", "totalAmount", "currency", "version", "createdAt", "updatedAt"},
	Filterable: []string{"projectId", "status", "year", "category", "fundingSource"},
	Sortable:   []string{"title", "totalAmount", "createdAt", "updatedAt"},
	Expandable: []string{"lineItems", "actuals"},
}

// BudgetLineItemODataFields defines allowed OData fields for line item queries
var BudgetLineItemODataFields = AllowedFields{
	Selectable: []string{"id", "budgetId", "category", "description", "budgetedAmount", "year", "fundingSource"},
	Filterable: []string{"budgetId", "category", "year", "fundingSource"},
	Sortable:   []string{"category", "budgetedAmount", "year"},
	Expandable: []string{"budget", "actuals"},
}

// BudgetActualODataFields defines allowed OData fields for actual queries
var BudgetActualODataFields = AllowedFields{
	Selectable: []string{"id", "lineItemId", "amount", "description", "recordedDate", "source", "externalRef"},
	Filterable: []string{"lineItemId", "recordedDate", "source"},
	Sortable:   []string{"amount", "recordedDate"},
	Expandable: []string{"lineItem"},
}

// IsFieldAllowed checks if a field is in the allowed list
func (af AllowedFields) IsSelectAllowed(field string) bool {
	for _, f := range af.Selectable {
		if f == field {
			return true
		}
	}
	return false
}

// IsFilterAllowed checks if a field can be filtered
func (af AllowedFields) IsFilterAllowed(field string) bool {
	for _, f := range af.Filterable {
		if f == field {
			return true
		}
	}
	return false
}

// IsSortAllowed checks if a field can be sorted
func (af AllowedFields) IsSortAllowed(field string) bool {
	for _, f := range af.Sortable {
		if f == field {
			return true
		}
	}
	return false
}

// IsExpandAllowed checks if a relation can be expanded
func (af AllowedFields) IsExpandAllowed(field string) bool {
	for _, f := range af.Expandable {
		if f == field {
			return true
		}
	}
	return false
}
