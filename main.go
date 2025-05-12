package main

import (
	"github.com/darnellsylvain/auth52/cmd/api"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	api := api.NewAPI()
	api.ListenAndServe()
}