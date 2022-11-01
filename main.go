package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maksimartemev/golang-db-pg-example/model"
	"github.com/maksimartemev/golang-db-pg-example/store"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		f(w, r)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId != "" {
		userIdInt, _ := strconv.Atoi(userId)

		user, _ := store.GetById(userIdInt)
		if user.ID < 1 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)

	var regUser model.RegisterUser
	err := decoder.Decode(&regUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdUser, err := store.Create(regUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(createdUser)
}

func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(users)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", logging(HomeHandler))
	router.HandleFunc("/register", logging(RegisterHandler))
	router.HandleFunc("/list", logging(ListUsersHandler))

	log.Fatal(http.ListenAndServe(":8081", router))
}
