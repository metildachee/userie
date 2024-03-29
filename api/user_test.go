package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/logger"
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
	router.HandleFunc("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetUser(w, r)
	})
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.NotNil(t, err, "json decoder err")
	require.EqualValues(t, models.User{}, user, "response is nil")
	assert.EqualValues(t, http.StatusNotFound, resp.Code, "response code is not ok")
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
	limit, offset, users := 2, 1, make([]models.User, 0)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/limit=%d&offset=%d", limit, offset), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/users/limit={limit}&offset={offset}", func(w http.ResponseWriter, r *http.Request) {
		GetAll(w, r)
	})
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&users)
	assert.Nil(t, err, "json decoder err")
	assert.EqualValues(t, limit, len(users), "does not conform to limit")
	assert.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")
}

func TestGetAllDefault(t *testing.T) {
	limit, users := 10, make([]models.User, 0)
	req, _ := http.NewRequest(http.MethodGet, "/api/users", nil)
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
	setup()
	userId, user := "1", models.User{}

	// can get
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%s", userId), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", GetUser)
	router.ServeHTTP(resp, req)

	err := json.NewDecoder(resp.Body).Decode(&user)
	require.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")
	assert.Nil(t, err, "json decoder err")
	require.NotEqualValues(t, models.User{}, user, "response is nil")

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
	router.HandleFunc("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetUser(w, r)
	})
	router.ServeHTTP(resp, req)

	u := models.User{}
	err = json.NewDecoder(resp.Body).Decode(&u)
	assert.NotNil(t, err, "json decoder err")
	require.EqualValues(t, models.User{}, u, "response should be nil")
	assert.EqualValues(t, http.StatusNotFound, resp.Code, "response code is not ok")
}

func TestCreateUser(t *testing.T) {
	setup()
	u := models.User{
		Name:        "metchee",
		DOB:         int32(time.Now().AddDate(-1, -1, -30).Unix()),
		Address:     "Kent Ridge",
		Description: "default user info",
		Ctime:       int32(time.Now().Unix()),
	}

	jsonBody, err := json.Marshal(u)
	require.Nil(t, err, "should not have error when marshal")

	fmt.Println(string(jsonBody))

	req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(jsonBody))
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/users", CreateUser)
	router.ServeHTTP(resp, req)

	var id int
	err = json.NewDecoder(resp.Body).Decode(&id)
	fmt.Println(id)
	require.Nil(t, err, "json decoder err=", err)
	require.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")
	assert.NotNil(t, id, "id should be auto incremented")
}

func TestUpdateUser(t *testing.T) {
	var (
		newName = "meow meow"
	)
	userId, user := "4", models.User{}

	// get
	getReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%s", userId), nil)
	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", GetUser)
	router.ServeHTTP(resp, getReq)

	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.Nil(t, err, "json decoder err")
	require.NotEqualValues(t, newName, user.Name, "name is already the same")
	assert.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")

	// update
	user.Name = newName
	jsonBody, err := json.Marshal(user)
	require.Nil(t, err, "should not have error when marshal")
	updateReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/user/%s", userId), bytes.NewBuffer(jsonBody))
	resp = httptest.NewRecorder()
	fmt.Println(string(jsonBody))

	router.HandleFunc("/api/user/{id}", UpdateUser)
	router.ServeHTTP(resp, updateReq)

	require.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")

	// get again
	router.ServeHTTP(resp, getReq)
	require.EqualValues(t, http.StatusOK, resp.Code, "response code is not ok")
	u := models.User{}
	err = json.NewDecoder(resp.Body).Decode(&u)
	assert.Nil(t, err, "json decoder err")
	require.EqualValues(t, newName, user.Name, "name is already the same")
}

func setup() {
	lf, err := os.OpenFile("user_server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("failed to open log file: %v", err)
	}
	defer lf.Close()
	defer logger.Init("info logger", true, true, lf).Close()
}
