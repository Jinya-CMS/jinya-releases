package utils

type ErrorDetails struct {
	EntityType string `json:"entityType"`
	ErrorType  string `json:"errorType"`
	Message    string `json:"message"`
}
