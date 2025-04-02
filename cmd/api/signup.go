package api

import (
	"fmt"
	"net/http"

	validator "github.com/darnellsylvain/auth52/internal"
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
	err := readJSON(w, r, params)
	if err != nil {
		badRequestError("Could not read signup params %v", err)

		return
	}
	
	v := validator.New()

	// Do validation checks on the params passed in. Makesure email and password in right format
	// Ensure both fields are not empty
	if ValidateSignupParams(v, params); !v.Valid() {
		handleError(validationError(v.Errors), w)
		fmt.Println(v.Errors)
		return
	}
	// Ensure Email is correctly formatted 
	// Ensure password is correctly formatted

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
	// sendJSON(w, http.StatusOK, user)
}

