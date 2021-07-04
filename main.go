package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metildachee/userie/api"
)

var (
	port string
)

func main() {
	flag.StringVar(&port, "port", ":8080", "server listen port")
	flag.Parse()

	r := mux.NewRouter()
	prefix := r.PathPrefix("/api").Subrouter()

	us := prefix.PathPrefix("/users").Subrouter()
	us.HandleFunc("/limit={limit}", api.GetAll).Methods(http.MethodGet)
	us.HandleFunc("/", api.CreateUser).Methods(http.MethodPost)
	us.HandleFunc("/{id}", api.UpdateUser).Methods(http.MethodPut)
	us.HandleFunc("/{id}", api.DeleteUser).Methods(http.MethodDelete)

	u := prefix.Path("/user").Subrouter()
	u.HandleFunc("/{id}", api.GetUser).Methods(http.MethodGet)
	
	log.Fatal(http.ListenAndServe(port, r))
}
