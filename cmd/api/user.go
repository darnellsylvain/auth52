package api

import (
	"net/http"

	"github.com/darnellsylvain/auth52/models"
)


func (api *API) GetUser(w http.ResponseWriter, r *http.Request) {
	// Middleware should check if this user is able to access this resource
	// Verify user
	// Get the user ID from req


	// Access user modela and find user by the ID
	// Return the user in JSON respon

	sendJSON(w, http.StatusOK, &models.User{
		Email:     "alice@example.com",

	})
}
