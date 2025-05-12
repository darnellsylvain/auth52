package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/darnellsylvain/auth52/internal/auth"
	validator "github.com/darnellsylvain/auth52/internal/validator"
	"github.com/darnellsylvain/auth52/models"
	"github.com/jackc/pgx/v5"
)


type LoginParams struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type LoginResponse struct {
	User *models.User 	`json:"user"`
	Token string		`json:"token"`
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

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
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

	
	token, err := auth.MakeJWT(user.ID, user.Email)
	if err != nil {
		api.serverErrorResponse(w, r, err)
	}

	response := &LoginResponse{
		User: 	user,
		Token: 	token,
	}

	sendJSON(w, http.StatusOK, response, nil)
}