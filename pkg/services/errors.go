package services

import "encoding/json"

type DomainError struct {
	statusCode int
	message    string
	details    map[string]interface{}
}

func NewDomainError(statusCode int, message string, details map[string]interface{}) error {
	return &DomainError{
		statusCode: statusCode,
		message:    message,
		details:    details,
	}
}

func (d *DomainError) Error() string {
	errorJson, _ := json.Marshal(map[string]interface{}{
		"status":  "error",
		"code":    d.statusCode,
		"message": d.message,
		"details": d.details,
	})
	return string(errorJson)
}

func (d *DomainError) StatusCode() int {
	return d.statusCode
}
