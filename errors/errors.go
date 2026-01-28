package errors

import "errors"

// Custom error types for consistent error handling across the application
var (
	// ErrNotFound is returned when a requested resource does not exist
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput is returned when the input data is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrConflict is returned when there's a conflict (e.g., duplicate entry)
	ErrConflict = errors.New("resource already exists")

	// ErrInternal is returned when an internal server error occurs
	ErrInternal = errors.New("internal server error")
)
