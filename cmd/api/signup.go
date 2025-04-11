package api

import (
	"context"
	"errors"
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

var ErrDuplicateEmail = errors.New("email already in use")

type SignupParams struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

func (api *API) Signup(w http.ResponseWriter, r *http.Request) {
	params := &SignupParams{}
	
	err := readJSON(w, r, params)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}
	
	v := validator.New()
	if ValidateSignupParams(v, params); !v.Valid() {
		api.failedValidationResponse(w, r, v.Errors)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	tx, err := api.db.Begin(ctx)
	if err != nil {
		api.serverErrorResponse(w, r, err)
		return
	}
	defer tx.Rollback(ctx)

	user, err := signupNewUser(ctx, tx, *params)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
    	api.serverErrorResponse(w, r, err)
    	return
	}

	sendJSON(w, http.StatusOK, user, nil)
}

func signupNewUser(ctx context.Context, tx pgx.Tx, params SignupParams) (*models.User, error) {
	var id pgtype.UUID
	err := tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, params.Email).Scan(&id)
	if err == nil {
		return nil, ErrDuplicateEmail
	} else  {
    	switch err {
    		case pgx.ErrNoRows:
        	break
    		default:
        		return nil, err
    	}
	}

	user, err := models.NewUser(params.Email, params.Password)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, query, user.Name, user.Email, user.EncryptedPassword, user.Activated, user.Provider)
	if err != nil {
		return nil, err
	}

	return user, nil
}