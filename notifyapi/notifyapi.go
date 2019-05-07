package notifyapi

import (
	"net/http"
	"strings"
)

// Error contains an error response from the server.
type Error struct {
	// Code is the HTTP response status code and will always be populated.
	Code int `json:"status_code"`
	// Header contains the response header fields from the server.
	Header http.Header

	Errors []ErrorItem `json:"errors"`
}

type ErrorItem struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	var allErrors []string
	for _, v := range e.Errors {
		allErrors = append(allErrors, v.Message)
	}
	return strings.Join(allErrors, ", ")
}
