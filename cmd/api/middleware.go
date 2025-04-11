package api

import (
	"log"
	"net/http"
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
