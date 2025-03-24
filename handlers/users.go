package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Main user interace
type UserService struct {
	l *log.Logger
}

func NewUser(l *log.Logger) *UserService {
	return &UserService{l}
}

func (u *UserService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		u.GetUsers(w, r)
		return
	case http.MethodPost:
		u.CreateUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}


}

func (u *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	// var input struct {
	// 	Name     string `json:"name"`
	// 	Email    string `json:"email"`
	// 	Password string `json:"password"`
	// }
	
	user := &data.User{}
	err := user.FromJSON(r.Body)

	if err != nil {
		http.Error(w, "cannot decode body", http.StatusBadRequest)
	}

	// user.SetPassword(input.Password)
	if err != nil {
		fmt.Println("something went wrong")
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(user)

}

func (u *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	lu := data.GetUsers()
	err := lu.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to Marshal json", http.StatusInternalServerError)
	}

}