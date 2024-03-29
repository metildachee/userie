package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/metildachee/userie/api"
	"github.com/metildachee/userie/utilities"
	"github.com/opentracing/opentracing-go"
)

var (
	logFilePath string
)

func main() {
	configFilePath := flag.String("configFilePath", "configuration.yml", "configuration file path")
	logFilePath := flag.String("logFilePath", "user_server.log", "user server info file path")
	verbose := flag.Bool("verbose", true, "some boolean")
	flag.Parse()

	// Fetch and set config
	if *configFilePath == "" || *logFilePath == "" {
		log.Fatal("invalid file path")
	}

	env, err := utilities.SetConfig(*configFilePath)
	if err != nil {
		log.Fatal("configuration cannot be read, aborting boot up")
	}

	// Init logging
	lf, err := os.OpenFile(*logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("failed to open log file: %v", err)
	}
	defer lf.Close()
	defer logger.Init("info logger", *verbose, *verbose, lf).Close()

	// Init tracer
	tracer, closer := utilities.InitJaeger(env.GetServiceName())
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	opentracing.ContextWithSpan(context.Background(), tracer.StartSpan("service started"))

	// Init http
	r := mux.NewRouter()
	prefix := r.PathPrefix("/api").Subrouter()

	u := prefix.PathPrefix("/user").Subrouter()
	u.HandleFunc("/{id}", api.GetUser).Methods(http.MethodGet)
	u.HandleFunc("", api.UpdateUser).Methods(http.MethodPut)
	u.HandleFunc("/{id}", api.DeleteUser).Methods(http.MethodDelete)
	u.HandleFunc("", api.CreateUser).Methods(http.MethodPost)

	us := prefix.PathPrefix("/users").Subrouter()
	us.HandleFunc("/limit={limit}&offset={offset}", api.GetAll).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(env.GetServerEndpoint(), r))
}
