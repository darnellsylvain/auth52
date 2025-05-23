package api

import (
	"errors"
	"net/http"

	"github.com/darnellsylvain/auth52/internal/auth"
	"github.com/darnellsylvain/auth52/models"
)

func (api *API) GetUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		api.unauthorizedResponse(w, r, errors.New("unauthorized"))
		return
	}

	sendJSON(w, http.StatusOK, &models.User{
		Email: claims.Email,
	}, nil)
}
