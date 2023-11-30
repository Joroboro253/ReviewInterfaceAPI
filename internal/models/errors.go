package models

import "net/http"

type APIError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e *APIError) Error() string {
	return e.Title
}

var (
	ErrReviewNotFound  = &APIError{Status: http.StatusNotFound, Title: "Review Not Found", Detail: "The requested review does not exist"}
	ErrInvalidInput    = &APIError{Status: http.StatusBadRequest, Title: "Invalid Input", Detail: "The provided input is not valid"}
	ErrDatabaseProblem = &APIError{Status: http.StatusInternalServerError, Title: "Database Problem", Detail: "A problem occurred with the database"}
	ErrInternal        = &APIError{Status: http.StatusInternalServerError, Title: "Internal Error", Detail: "An internal error occurred"}
)
