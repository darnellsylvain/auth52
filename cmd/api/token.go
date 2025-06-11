package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/darnellsylvain/auth52/internal/auth"
	"github.com/darnellsylvain/auth52/internal/database"
)

type RefreshParams struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken       string    `json:"accessToken"`
	AccessTokenExpiry time.Time `json:"accessTokenExpiry"`
	RefreshToken      string    `json:"refreshToken"`
}

func (api *API) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := api.db.Begin(ctx)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback(ctx)

	qtx := api.queries.WithTx(tx)

	user, err := qtx.GetUserByValidSession(ctx, refreshToken)
	if err != nil {
		api.unauthorizedResponse(w, r, errors.New("invalid refresh token"))
		return
	}

	newAccessToken, newClaims, err := auth.CreateToken(user.UserID, user.Email)
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}

	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}

	_, err = qtx.UpdateSessionByToken(ctx, database.UpdateSessionByTokenParams{
		RefreshToken:   refreshToken,
		RefreshToken_2: newRefreshToken,
		ExpiresAt:      time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}

	tx.Commit(ctx)

	response := &RefreshResponse{
		AccessToken:       newAccessToken,
		AccessTokenExpiry: newClaims.ExpiresAt.Time,
		RefreshToken:      newRefreshToken,
	}

	sendJSON(w, http.StatusOK, response, nil)
}
