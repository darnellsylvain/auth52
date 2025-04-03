package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type HTTPError struct {
	Code            int 				`json:"code"`
	Message         string				`json:"message"`
	InternalError   error				`json:"-"`
	InternalMessage string				`json:"-"`
	Errors			map[string]string	`json:"errors"`
}

func (e *HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func badRequestError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusBadRequest, fmtString, args...)
}

func serverError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusInternalServerError, fmtString, args...)
}

func validationError(errors map[string]string) *HTTPError {
		return &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "validation failed",
		Errors:  errors,
	}
}


func unprocessableEntityError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnprocessableEntity, fmtString, args...)
}


func httpError(code int, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: fmt.Sprintf(fmtString, args...),
	}
}

func handleError(err *HTTPError, w http.ResponseWriter) {
	sendJSON(w, err.Code, err)
}

func isDuplicateError(err error) bool {
    var pgErr *pgconn.PgError
    return errors.As(err, &pgErr) && pgErr.Code == "23505"
}