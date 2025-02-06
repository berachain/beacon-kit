package handlers

// HTTPError represents an HTTP error response.
type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return e.Message
}

// StatusCode returns the HTTP status code
func (e *HTTPError) StatusCode() int {
	return e.Code
}
