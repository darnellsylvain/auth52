package api

import (
	"net/http"

	"github.com/darnellsylvain/auth52/models"
)

type SignupParams struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

func (api *API) Signup(w http.ResponseWriter, r *http.Request) {
	// Get params passed in
	params := &SignupParams{}
	// read params/decode
	readJSON(w, r, params)
	
	// Do validation checks on the params passed in. Makesure email and password in right format

	// Do a DB query to find the email, if it exists send an error

	// If a user doesn't exist already, create a new user and store in DB. This can be own function
		// Create new model instance
	user, err := models.NewUser(params.Email, params.Password)
	if err != nil {
		// handle error
		sendJSON(w, http.StatusInternalServerError, nil)
	}
		// Set any fields you need
	user.Provider = "email"
		// Store in DB
		// Handle any errors


	// Send a response back with HTTP Ok and user
	sendJSON(w, http.StatusOK, user)

}