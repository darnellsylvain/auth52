package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/darnellsylvain/auth52/internal/database"
	validator "github.com/darnellsylvain/auth52/internal/validator"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/darnellsylvain/auth52/models"
)

var ErrDuplicateEmail = errors.New("email already in use")

type SignupParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := signupNewUser(ctx, api.queries, *params)
	if err != nil {
		api.badRequestError(w, r, err)
		return
	}

	sendJSON(w, http.StatusOK, user, nil)
}

func signupNewUser(ctx context.Context, queries *database.Queries, params SignupParams) (*models.User, error) {
	user, err := models.NewUser(params.Email, params.Password)
	if err != nil {
		return nil, err
	}

	createdUser, err := queries.CreateUser(ctx, database.CreateUserParams{
		Name:              user.Name,
		Email:             user.Email,
		EncryptedPassword: user.EncryptedPassword,
		Activated:         user.Activated,
		Provider:          user.Provider,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}

	user.ID = createdUser.ID
	user.CreatedAt = createdUser.CreatedAt
	return user, nil
}
