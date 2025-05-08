package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/darnellsylvain/auth52/internal/database"
	"github.com/darnellsylvain/auth52/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type envelope map[string]any


func sendJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	
	return nil
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {	
		return err
	}

	return nil
}


func toTimeOrZero(tz pgtype.Timestamptz) time.Time {
    if tz.Valid {
        return tz.Time
    }
    return time.Time{}
}

func FromDBUser(u database.FindUserByEmailRow) *models.User {
	return &models.User{
		ID:                	u.ID,
		CreatedAt:			toTimeOrZero(u.CreatedAt),
		Name:              	u.Name,
		Email:             	u.Email,
		EncryptedPassword: 	u.EncryptedPassword,
		Activated:         	u.Activated,
		Provider:         	u.Provider,
	}
}