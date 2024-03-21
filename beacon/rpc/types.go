package rpc

import "fmt"

const JsonMediaType = "application/json"

type HasStatusCode interface {
	StatusCode() int
}

// DefaultJsonError is a JSON representation of a simple error value, containing only a message and an error code.
type DefaultJsonError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *DefaultJsonError) StatusCode() int {
	return e.Code
}

func (e *DefaultJsonError) Error() string {
	return fmt.Sprintf("HTTP request unsuccessful (%d: %s)", e.Code, e.Message)
}
