package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/darnellsylvain/auth52/internal/auth"
)

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = api.queries.RevokeSessionByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.unauthorizedResponse(w, r, errors.New("invalid token"))

		} else {
			api.serverErrorResponse(w, r, err)

		}
		return

	}

	sendJSON(w, http.StatusOK, map[string]string{
		"status":  "200",
		"message": "logged out successfully",
	}, nil)

}
