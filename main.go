package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metildachee/userie/handler"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/users", handler.GetUsers).Methods("GET")
	r.HandleFunc("/api/user/{id}", handler.GetUser).Methods("GET")
	r.HandleFunc("/api/users", handler.CreateUser).Methods("CREATE")
	r.HandleFunc("/api/users/{id}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/users/{id}", handler.DeleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))
}
