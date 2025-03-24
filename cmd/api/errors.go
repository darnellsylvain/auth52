package api

import (
	"fmt"
	"net/http"
)



type HTTPError struct {
	Code				int
	Message 			string
	InternalError 		error
	InternalMessage 	string
}



func badRequestError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusBadRequest, fmtString, args...)
}

func httpError(code int, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		Code: 		code,
		Message: 	fmt.Sprintf(fmtString, args...),
	}
}