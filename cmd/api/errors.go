package api

import (
	"net/http"
)

type HttpError struct {
	Status  int
	Message string
	Errors  map[string]string
}

func (api *API) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	api.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (api *API) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message, "status": status}

	err := sendJSON(w, status, env, nil)
	if err != nil {
		api.logError(r, err)
		w.WriteHeader(500)
	}
}

func (api *API) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	api.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (api *API) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	api.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	api.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (api *API) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	api.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
