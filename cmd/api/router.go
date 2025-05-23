package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (api *API) NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/healthcheck", api.HealthCheck).Methods("GET")
	r.HandleFunc("/signup", api.Signup).Methods("POST")
	r.HandleFunc("/login", api.Login).Methods("GET")

	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.Use(api.RequireAuthorization)
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
