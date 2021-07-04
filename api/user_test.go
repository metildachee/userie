package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/metildachee/userie/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserInvalid(t *testing.T) {
	userId, user := 0, models.User{}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%d", userId), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", GetUser)
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.NotNil(t, err, "json decoder err")
	require.EqualValues(t, models.User{}, user, "response is nil")
	assert.EqualValues(t, http.StatusBadRequest, resp.Code, "response code is not ok")
}

func TestGetUserValid(t *testing.T) {
	userId, user := 1, models.User{}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%d", userId), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", GetUser)
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.Nil(t, err, "json decoder err")
	require.NotEqualValues(t, models.User{}, user, "response is nil")
	assert.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")
}

func TestGetAllWithLimit(t *testing.T) {
	limit, users := 2, make([]models.User, 0)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/limit=%d", limit), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/users/limit={limit}", GetAll)
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&users)
	assert.Nil(t, err, "json decoder err")
	assert.EqualValues(t, limit, len(users), "does not conform to limit")
	assert.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")
}

func TestGetAllDefault(t *testing.T) {
	limit, users := 10, make([]models.User, 0)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/users"), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/users", GetAll)
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&users)
	assert.Nil(t, err, "json decoder err")
	assert.LessOrEqual(t, len(users), limit, "should be lesser than or equal to limit")
	assert.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")
}

func TestDeleteUser(t *testing.T) {
	userId, user := "7", models.User{}

	// can get
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%s", userId), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", GetUser)
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.Nil(t, err, "json decoder err")
	require.NotEqualValues(t, models.User{}, user, "response is nil")
	assert.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")

	// delete
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/user/%s", userId), nil)
	resp = httptest.NewRecorder()

	router = mux.NewRouter()
	router.HandleFunc("/api/user/{id}", DeleteUser)
	router.ServeHTTP(resp, req)

	assert.EqualValues(t, http.StatusNoContent, resp.Code, "response code is not ok")

	// cannot get
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%s", userId), nil)
	resp = httptest.NewRecorder()

	router = mux.NewRouter()
	router.HandleFunc("/api/user/{id}", GetUser)
	router.ServeHTTP(resp, req)

	u := models.User{}
	err = json.NewDecoder(resp.Body).Decode(&u)
	assert.NotNil(t, err, "json decoder err")
	require.EqualValues(t, models.User{}, u, "response should be nil")
	assert.EqualValues(t, http.StatusNotFound, resp.Code, "response code is not ok")
}
