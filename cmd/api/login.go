package api

import (
	"context"
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
		badRequestError("Could not read signup params %v", err)
		return 
	}

	v := validator.New()
	if ValidateLoginParams(v, params); !v.Valid() {
		handleError(validationError(v.Errors), w)
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
			handleError(badRequestError("no user with this email %v", err), w)
			return 
		}
		return

	}

	ok := user.Authenticate(params.Password)
	if !ok {
		handleError(badRequestError("email or password is incorrect"), w)
		return 
	}

	sendJSON(w, http.StatusOK, user)
}