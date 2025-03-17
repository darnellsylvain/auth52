package main

import "github.com/darnellsylvain/auth52/cmd/api"


func main() {
	api := api.NewAPI()
	api.ListenAndServe()
}