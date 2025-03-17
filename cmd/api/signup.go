package api

import "net/http"

type SignupParams struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

func (api *API) Signup(w http.ResponseWriter, r *http.Request) {
	params := &SignupParams{}

	readJSON(w, r, params)
	sendJSON(w, http.StatusOK, params)

}