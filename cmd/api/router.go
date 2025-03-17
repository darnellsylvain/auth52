package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	mux *mux.Router
}

func NewRouter() *Router {
	return &Router{
		mux: mux.NewRouter(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
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
