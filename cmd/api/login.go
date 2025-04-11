package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	validator "github.com/darnellsylvain/auth52/internal"
	"github.com/darnellsylvain/auth52/models"
	"github.com/jackc/pgx/v5"
)

var loginQuery = `SELECT id, created_at, name, email, encrypted_password, activated, provider FROM users WHERE email = $1`

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

	user := &models.User{}
	err = api.db.QueryRow(ctx, loginQuery, params.Email).Scan(&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.EncryptedPassword,
		&user.Activated,
		&user.Provider,
	)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			api.badRequestError(w, r, err)
			return 
		}
		return

	}

	ok := user.Authenticate(params.Password)
	if !ok {
		api.badRequestError(w, r, errors.New("email or password is incorrect"))
		return 
	}

	sendJSON(w, http.StatusOK, user, nil)
}