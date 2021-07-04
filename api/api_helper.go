package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getParam(field string, r *http.Request) string {
	params := mux.Vars(r)
	return params[field]
}

func writeJsonHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	return w
}