package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metildachee/userie/api"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/users/limit={limit}", api.GetAll).Methods("GET")
	r.HandleFunc("/api/user/{id}", api.GetUser).Methods("GET")
	r.HandleFunc("/api/users", api.CreateUser).Methods("POST")
	r.HandleFunc("/api/users/{id}", api.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", api.DeleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
