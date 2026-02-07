package ldn

import "errors"

// Validation errors for Activity.
var (
	ErrMissingContext = errors.New("@context is required")
	ErrMissingID      = errors.New("id is required")
	ErrMissingType    = errors.New("type is required")
	ErrMissingOrigin  = errors.New("origin with id is required")
	ErrMissingTarget  = errors.New("target with id and inbox is required")
	ErrMissingObject  = errors.New("object with id is required")
)
