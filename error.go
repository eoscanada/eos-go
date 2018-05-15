package eos

import "errors"

// APIError represents the errors as reported by the server
type APIError struct {
	Code        int
	Message     string
	ErrorStruct struct {
		Code    int
		Name    string
		What    string
		Details []struct {
			Message    string
			File       string
			LineNumber int `json:"line_number"`
			Method     string
		}
	} `json:"error"`
}

func (e APIError) Error() error {
	return errors.New(e.String())
}

func (e APIError) String() string {
	return e.Message
}
