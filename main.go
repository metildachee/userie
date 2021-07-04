package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metildachee/userie/api"
	"github.com/metildachee/userie/utilities"
	"github.com/opentracing/opentracing-go"
)

var (
	configFilePath string
)

func main() {
	flag.StringVar(&configFilePath, "configFilePath", "configuration.yml", "configuration file path")
	flag.Parse()

	if configFilePath == "" {
		log.Fatal("invalid configuration file path")
	}

	env, err := utilities.SetConfig(configFilePath)
	if err != nil {
		log.Fatal("configuration cannot be read, aborting boot up")
	}

	tracer, closer := utilities.InitJaeger(env.GetServiceEnvName())
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	ctx := opentracing.ContextWithSpan(context.Background(), tracer.StartSpan("service started"))

	r := mux.NewRouter()
	prefix := r.PathPrefix("/api").Subrouter()

	us := prefix.PathPrefix("/users").Subrouter()
	us.HandleFunc("/limit={limit}", func(w http.ResponseWriter, r *http.Request) {
		api.GetAll(w, r, ctx)
	}).Methods(http.MethodGet)
	us.HandleFunc("/", api.CreateUser).Methods(http.MethodPost)
	us.HandleFunc("/{id}", api.UpdateUser).Methods(http.MethodPut)
	us.HandleFunc("/{id}", api.DeleteUser).Methods(http.MethodDelete)

	u := prefix.Path("/user").Subrouter()
	u.HandleFunc("/{id}", api.GetUser).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(env.GetServerEndpoint(), r))
}
