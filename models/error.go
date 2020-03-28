package models

import "net/http"

// JSONError is the struct representing a JSON formatted error
type JSONError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// NewBadRequestError return an error 400 JSON formatted
func NewBadRequestError() *JSONError {
	return &JSONError{
		Status: http.StatusBadRequest,
		Error:  "Bad Request Error",
	}
}

// NewUnauthorizedError return an error 401 JSON formatted
func NewUnauthorizedError() *JSONError {
	return &JSONError{
		Status: http.StatusUnauthorized,
		Error:  "Unauthorized",
	}
}

// NewForbiddenError return an error 403 JSON formatted
func NewForbiddenError() *JSONError {
	return &JSONError{
		Status: http.StatusForbidden,
		Error:  "Unauthorized",
	}
}

// NewInternalServerError returns an error 500 JSON formatted
func NewInternalServerError() *JSONError {
	return &JSONError{
		Status: http.StatusInternalServerError,
		Error:  "Internal Server Error",
	}
}
