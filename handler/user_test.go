package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/metildachee/userie/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserInvalid(t *testing.T) {
	userId, user := 0, model.User{}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%d", userId), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", GetUser)
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.NotNil(t, err, "json decoder err")
	require.EqualValues(t, model.User{}, user, "response is nil")
	assert.EqualValues(t, http.StatusBadRequest, resp.Code, "response code is not ok")
}

func TestGetUserValid(t *testing.T) {
	userId, user := 1, model.User{}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%d", userId), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", GetUser)
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.Nil(t, err, "json decoder err")
	require.NotEqualValues(t, model.User{}, user, "response is nil")
	assert.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")
}
