package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/metildachee/userie/api"
	"github.com/metildachee/userie/utilities"
)

var (
	configFilePath string
)

func main() {
	flag.StringVar(&configFilePath, "configFilePath", "configuration.json", "configuration file path")
	flag.Parse()

	if configFilePath == "" {
		log.Fatal("invalid configuration file path")
	}

	env, err := utilities.SetConfig(configFilePath)
	if err != nil {
		log.Fatal("configuration cannot be read, aborting boot up")
	}

	r := mux.NewRouter()
	prefix := r.PathPrefix("/api").Subrouter()

	us := prefix.PathPrefix("/users").Subrouter()
	us.HandleFunc("/limit={limit}", api.GetAll).Methods(http.MethodGet)
	us.HandleFunc("/", api.CreateUser).Methods(http.MethodPost)
	us.HandleFunc("/{id}", api.UpdateUser).Methods(http.MethodPut)
	us.HandleFunc("/{id}", api.DeleteUser).Methods(http.MethodDelete)

	u := prefix.Path("/user").Subrouter()
	u.HandleFunc("/{id}", api.GetUser).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(os.Getenv(env.ServerPort), r))
}
