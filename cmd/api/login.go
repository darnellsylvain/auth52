package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/darnellsylvain/auth52/internal/auth"
	"github.com/darnellsylvain/auth52/internal/database"
	validator "github.com/darnellsylvain/auth52/internal/validator"
	"github.com/darnellsylvain/auth52/models"
	"github.com/jackc/pgx/v5"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User              *models.User `json:"user"`
	AccessToken       string       `json:"accessToken"`
	AccessTokenExpiry time.Time    `json:"accessTokenExpiry"`
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	params := &LoginParams{}

	if err := readJSON(w, r, params); err != nil {
		api.badRequestError(w, r, err)
		return
	}

	v := validator.New()
	if ValidateLoginParams(v, params); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userDB, err := api.queries.FindUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			api.badRequestError(w, r, errors.New("email or password is incorrect"))
			return
		}
		api.serverErrorResponse(w, r, err)
		return
	}

	user := models.FromDBUser(userDB)

	ok := user.Authenticate(params.Password)
	if !ok {
		api.badRequestError(w, r, errors.New("email or password is incorrect"))
		return
	}

	accessToken, accessClaims, err := auth.CreateToken(user.ID, user.Email)
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}

	ipAddr := GetIPAddressFromRequest(r)

	_, err = api.queries.CreateSession(ctx, database.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
		IpAddress:    ipAddr,
	})
	if err != nil {
		api.serverErrorResponse(w, r, errors.New("error creating session"))
		return
	}

	response := &LoginResponse{
		User:              user,
		AccessToken:       accessToken,
		AccessTokenExpiry: accessClaims.ExpiresAt.Time,
	}

	sendJSON(w, http.StatusOK, response, nil)
}
