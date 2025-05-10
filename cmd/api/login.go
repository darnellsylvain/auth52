package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	validator "github.com/darnellsylvain/auth52/internal/validator"
	"github.com/darnellsylvain/auth52/models"
	"github.com/jackc/pgx/v5"
)


type LoginParams struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	params := &LoginParams{}

	err := readJSON(w, r, params)
	if err != nil {
		api.badRequestError(w, r, err)
		return 
	}

	v := validator.New()
	if ValidateLoginParams(v, params); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return 
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	userDB, err := api.queries.FindUserByEmail(ctx, params.Email)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			api.badRequestError(w, r, errors.New("email or password is incorrect"))
			return 
		}
		return
	}

	user := models.FromDBUser(userDB)

	ok := user.Authenticate(params.Password)
	if !ok {
		api.badRequestError(w, r, errors.New("email or password is incorrect"))
		return 
	}

	sendJSON(w, http.StatusOK, user, nil)
}