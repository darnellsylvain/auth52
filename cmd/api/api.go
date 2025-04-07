package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darnellsylvain/auth52/storage"
)

type API struct {
	handler 	http.Handler
	db 			*storage.Connection
	// config
	version string
}


func NewAPI() *API {
	api := &API{
		version: "1",
	}

	// Initialise DB
	db, err := storage.Dial()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	api.db = db

	// Initiallize Router
	router := NewRouter()
	router.mux.HandleFunc("/healthcheck", api.HealthCheck)
	router.mux.HandleFunc("/user", api.GetUser)
	router.mux.HandleFunc("/signup", api.Signup).Methods("POST")
	router.mux.HandleFunc("/login", api.Login).Methods("GET")

	api.handler = router

	return api
}

func (api *API) ListenAndServe() {
	l := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	server := &http.Server{
		Addr: ":8080",
		Handler: api.RecoverPanic(api.handler),
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: time.Minute,
	}

	done := make(chan struct{})

	go func() {
		waitForTermination(l, done)
		ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
		defer cancel()
		api.Shutdown(ctx)
		server.Shutdown(ctx)
	}()

	err := server.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}
}

func waitForTermination(log *log.Logger, done <-chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	signal := <-sigChan
	log.Fatalf("Recieved shutdown signal from %s", signal)

	<-done
	log.Println("Shutting down...")
}

func (api *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, map[string]string{
		"version":     api.version,
		"name":        "Auth 52",
		"description": "Auth 52 is a user registration and authentication API",
	})
}


func (api *API) Shutdown(ctx context.Context) {
	if api.db != nil {
		api.db.Close()
		log.Println("Database connection closed")
	}
}