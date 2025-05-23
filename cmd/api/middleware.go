package api

import (
	"log"
	"net/http"

	"github.com/darnellsylvain/auth52/internal/auth"
)

func (api *API) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				w.Header().Set("Connection", "close")
				sendJSON(w, http.StatusInternalServerError, err, nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (api *API) RequireAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			api.unauthorizedResponse(w, r, err)
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			api.unauthorizedResponse(w, r, err)
			return
		}

		ctx := auth.SetClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
