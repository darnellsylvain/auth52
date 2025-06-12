package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (api *API) NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(api.RecoverPanic)
	r.Use(api.rateLimiter)
	apiRouter := r.PathPrefix("/api/" + api.version).Subrouter()

	apiRouter.HandleFunc("/healthcheck", api.HealthCheck).Methods("GET")
	apiRouter.HandleFunc("/signup", api.Signup).Methods("POST")
	apiRouter.HandleFunc("/login", api.Login).Methods("GET")
	apiRouter.HandleFunc("/logout", api.Logout).Methods("POST")
	apiRouter.HandleFunc("/refresh", api.Refresh).Methods("POST")

	userRouter := apiRouter.PathPrefix("/user").Subrouter()
	userRouter.Use(api.requireAuthorization)
	userRouter.HandleFunc("", api.GetUser).Methods("GET")

	return r

}

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func handler(fn apiHandler) http.HandlerFunc {
	return fn.serve
}

func (h apiHandler) serve(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		// Handle Error
	}
}
