package odata

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
)

// TODO: This file has to be moved somewhere else

// Parser implements OData query string parsing
type Parser struct{}

// NewParser creates a new OData query parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses OData query parameters into a structured query
func (p *Parser) Parse(queryParams url.Values) (entities.ODataQuery, error) {
	query := entities.ODataQuery{}

	// Parse $select
	if selectStr := queryParams.Get("$select"); selectStr != "" {
		query.Select = parseCommaSeparated(selectStr)
	}

	// Parse $filter
	if filterStr := queryParams.Get("$filter"); filterStr != "" {
		filter, err := parseFilter(filterStr)
		if err != nil {
			return query, err
		}
		query.Filter = filter
	}

	// Parse $orderby
	if orderByStr := queryParams.Get("$orderby"); orderByStr != "" {
		query.OrderBy = parseOrderBy(orderByStr)
	}

	// Parse $top
	if topStr := queryParams.Get("$top"); topStr != "" {
		top, err := strconv.Atoi(topStr)
		if err != nil {
			return query, &ParseError{Field: "$top", Reason: "must be an integer"}
		}
		if top < 0 || top > 1000 {
			return query, &ParseError{Field: "$top", Reason: "must be between 0 and 1000"}
		}
		query.Top = &top
	}

	// Parse $skip
	if skipStr := queryParams.Get("$skip"); skipStr != "" {
		skip, err := strconv.Atoi(skipStr)
		if err != nil {
			return query, &ParseError{Field: "$skip", Reason: "must be an integer"}
		}
		if skip < 0 {
			return query, &ParseError{Field: "$skip", Reason: "must be non-negative"}
		}
		query.Skip = &skip
	}

	// Parse $expand
	if expandStr := queryParams.Get("$expand"); expandStr != "" {
		query.Expand = parseCommaSeparated(expandStr)
	}

	// Parse $count
	if countStr := queryParams.Get("$count"); countStr == "true" {
		query.Count = true
	}

	return query, nil
}

func parseCommaSeparated(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseOrderBy(s string) []entities.ODataOrderBy {
	parts := strings.Split(s, ",")
	result := make([]entities.ODataOrderBy, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		orderBy := entities.ODataOrderBy{}
		tokens := strings.Fields(trimmed)
		if len(tokens) >= 1 {
			orderBy.Field = tokens[0]
		}
		if len(tokens) >= 2 && strings.ToLower(tokens[1]) == "desc" {
			orderBy.Desc = true
		}
		result = append(result, orderBy)
	}
	return result
}

// parseFilter parses a simple OData filter expression
// Supports: field eq value, field ne value, field gt value, etc.
// Also supports 'and' and 'or' for simple cases
func parseFilter(s string) (*entities.ODataFilter, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	// Check for 'and'/'or' at top level (simple split)
	lowerS := strings.ToLower(s)
	if idx := strings.Index(lowerS, " and "); idx > 0 {
		left, err := parseFilter(s[:idx])
		if err != nil {
			return nil, err
		}
		right, err := parseFilter(s[idx+5:])
		if err != nil {
			return nil, err
		}
		return &entities.ODataFilter{
			Field:    left.Field,
			Operator: left.Operator,
			Value:    left.Value,
			And:      right,
		}, nil
	}

	if idx := strings.Index(lowerS, " or "); idx > 0 {
		left, err := parseFilter(s[:idx])
		if err != nil {
			return nil, err
		}
		right, err := parseFilter(s[idx+4:])
		if err != nil {
			return nil, err
		}
		return &entities.ODataFilter{
			Field:    left.Field,
			Operator: left.Operator,
			Value:    left.Value,
			Or:       right,
		}, nil
	}

	// Parse single comparison: field operator value
	operators := []string{" eq ", " ne ", " gt ", " ge ", " lt ", " le "}
	for _, op := range operators {
		if idx := strings.Index(lowerS, op); idx > 0 {
			field := strings.TrimSpace(s[:idx])
			value := strings.TrimSpace(s[idx+len(op):])

			// Remove quotes from string values
			if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
				value = value[1 : len(value)-1]
			}

			return &entities.ODataFilter{
				Field:    field,
				Operator: entities.ODataOperator(strings.TrimSpace(op)),
				Value:    parseValue(value),
			}, nil
		}
	}

	// Check for contains(field, 'value')
	if strings.HasPrefix(lowerS, "contains(") {
		inner := s[9 : len(s)-1] // Remove "contains(" and ")"
		parts := strings.SplitN(inner, ",", 2)
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes
			if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
				value = value[1 : len(value)-1]
			}
			return &entities.ODataFilter{
				Field:    field,
				Operator: entities.ODataOpContains,
				Value:    value,
			}, nil
		}
	}

	return nil, &ParseError{Field: "$filter", Reason: "unsupported filter expression: " + s}
}

func parseValue(s string) any {
	// Try to parse as integer
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	// Try to parse as float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	// Try to parse as boolean
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	// Return as string
	return s
}

// ParseError indicates a parsing error
type ParseError struct {
	Field  string
	Reason string
}

func (e *ParseError) Error() string {
	return "parse error for " + e.Field + ": " + e.Reason
}
