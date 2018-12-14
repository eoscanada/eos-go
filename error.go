package eos

import (
	"fmt"

	"github.com/eoscanada/eos-go/eoserr"
)

// APIError represents the errors as reported by the server
type APIError struct {
	Code        int    `json:"code"` // http code
	Message     string `json:"message"`
	ErrorStruct struct {
		Code    int              `json:"code"` // https://docs.google.com/spreadsheets/d/1uHeNDLnCVygqYK-V01CFANuxUwgRkNkrmeLm9MLqu9c/edit#gid=0
		Name    string           `json:"name"`
		What    string           `json:"what"`
		Details []APIErrorDetail `json:"details"`
	} `json:"error"`
}

func NewAPIError(httpCode int, msg string, e eoserr.Error) *APIError {
	newError := &APIError{
		Code:    httpCode,
		Message: msg,
	}
	newError.ErrorStruct.Code = e.Code
	newError.ErrorStruct.Name = e.Name
	newError.ErrorStruct.What = msg
	newError.ErrorStruct.Details = []APIErrorDetail{
		APIErrorDetail{
			File:       "",
			LineNumber: 0,
			Message:    msg,
			Method:     e.Name,
		},
	}

	return newError
}

type APIErrorDetail struct {
	Message    string `json:"message"`
	File       string `json:"file"`
	LineNumber int    `json:"line_number"`
	Method     string `json:"method"`
}

func (e APIError) Error() string {
	msg := e.Message
	msg = fmt.Sprintf("%s: %s", msg, e.ErrorStruct.What)

	for _, detail := range e.ErrorStruct.Details {
		msg = fmt.Sprintf("%s: %s", msg, detail.Message)
	}

	return msg
}
