package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/metildachee/userie/dao"
	"github.com/metildachee/userie/logger"
	"github.com/metildachee/userie/model"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := dao.GetUsers()
	if err != nil {
		logger.Print(fmt.Sprintf("get users from dao err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		logger.Print(fmt.Sprintf("get users json encoder err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if params["id"] == "" {
		logger.Print(fmt.Sprintf("missing params of id, params=%s", params), ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(params["id"])
	if err != nil || userId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := dao.GetUser(int32(userId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		logger.Print(fmt.Sprintf("get user json encoder err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	newUser := model.User{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		logger.Print(fmt.Sprintf("create users json decode err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := newUser.Validate(); err != nil {
		logger.Print(fmt.Sprintf("create users validation err=%s", err), ERROR)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := dao.CreateUser(newUser); err != nil {
		logger.Print(fmt.Sprintf("create users from dao err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	updatedUser := model.User{}
	if params["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(params["id"])
	if err != nil || userId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		logger.Print(fmt.Sprintf("update users json decode err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := updatedUser.Validate(); err != nil {
		logger.Print(fmt.Sprintf("update users invalid user err=%s", err), ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dao.UpdateUser(int32(userId), updatedUser); err != nil {
		logger.Print(fmt.Sprintf("update users dao err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if params["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(params["id"])
	if err != nil || userId <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dao.DeleteUser(int32(userId)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
