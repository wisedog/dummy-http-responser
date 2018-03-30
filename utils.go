package main

import (
	"crypto/rand"
	"fmt"
)

type errorResponse struct {
	ErrorType string `json:"error"`
	ErrorMsg  string `json:"error_msg"`
}

func (e *errorResponse) Error() string {
	return e.ErrorType
}

var (
	errorInvalidStatus = &errorResponse{"InvalidStatus", "Make sure that your dummy-status is an integer value"}
	errorNotFound      = &errorResponse{"NotFound", "Check your URL again"}
	errorInvalidID     = &errorResponse{"InvalidID", "Invalid ID. Check your URL again"}
	errorInvalidData   = &errorResponse{"InvalidData", "Invalid data"}
)

func getRandomString() string {
	n := 16
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
