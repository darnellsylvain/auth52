package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	validator "github.com/darnellsylvain/auth52/internal"
	"github.com/darnellsylvain/auth52/models"
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


	var id pgtype.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, params.Email).Scan(&id)
	if err != nil {
		fmt.Println(err)
		if isDuplicateError(err) {
			handleError(badRequestError("user already exists"), w)
		}
	}

	if id.Valid {
		v.AddError("user", "already exists")
		handleError(badRequestError("user already exists"), w)
		return
	}

	user, err := models.NewUser(params.Email, params.Password)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, nil)
		return
	}
	user.Provider = "email"

	_, err = tx.Exec(ctx, query, user.Name, user.Email, user.EncryptedPassword, user.Activated, user.Provider)
	if err != nil {
		if err.Error() == `duplicate key value violates unique constraint "users_email_key"` {
			handleError(badRequestError("user already exists"), w)
		}
		handleError(serverError("Error adding new user"), w)
		return
	}

	tx.Commit(ctx)

	sendJSON(w, http.StatusOK, user)
}

