package api

import (
	"context"
	"net/http"
	"time"

	validator "github.com/darnellsylvain/auth52/internal"
	"github.com/darnellsylvain/auth52/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)
var query = `INSERT INTO users (name, email, encrypted_password, activated, provider) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`

type SignupParams struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

func (api *API) Signup(w http.ResponseWriter, r *http.Request) {
	params := &SignupParams{}
	
	err := readJSON(w, r, params)
	if err != nil {
		badRequestError("Could not read signup params %v", err)
		return
	}
	
	v := validator.New()
	if ValidateSignupParams(v, params); !v.Valid() {
		handleError(validationError(v.Errors), w)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	tx, err := api.db.Begin(ctx)
	if err != nil {
		handleError(serverError("Error creating transactions"), w)
	}
	defer tx.Rollback(ctx)

	user, err := signupNewUser(ctx, tx, *params)
	if err != nil {
		handleError(serverError("Error creating user"), w)
	}

	tx.Commit(ctx)

	sendJSON(w, http.StatusOK, user)
}

func signupNewUser(ctx context.Context, tx pgx.Tx, params SignupParams) (*models.User, error) {
	var id pgtype.UUID
	err := tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, params.Email).Scan(&id)
	if err != nil {
		return nil, err
	}

	if id.Valid {
		return nil, badRequestError("user already exists")
	}

	user, err := models.NewUser(params.Email, params.Password)
	if err != nil {
		return nil, serverError("could not create user")
	}

	_, err = tx.Exec(ctx, query, user.Name, user.Email, user.EncryptedPassword, user.Activated, user.Provider)
	if err != nil {
		return nil, serverError("Error adding new user")
	}

	return user, nil
}