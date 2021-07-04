package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/metildachee/userie/dao/elasticsearch"
	"github.com/metildachee/userie/logger"
	"github.com/metildachee/userie/models"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	w = writeJsonHeader(w)

	dao, err := elasticsearch.NewDao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	limit := 10
	if queryLimit := getParam("limit", r); queryLimit != "" {
		if limit, err = strconv.Atoi(queryLimit); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	users, err := dao.GetAll(limit)
	if err != nil {
		logger.Print(fmt.Sprintf("get users from dao err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		logger.Print(fmt.Sprintf("get users json encoder err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w = writeJsonHeader(w)
	dao, err := elasticsearch.NewDao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := ""
	if userId = getParam("id", r); userId == "" {
		logger.Print(fmt.Sprintf("missing params of id, params=%s", userId), ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	if user, err = dao.GetById(userId); err != nil {
		if err.Error() == "nil hit" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		logger.Print(fmt.Sprintf("get user json encoder err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	dao, err := elasticsearch.NewDao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		logger.Print(fmt.Sprintf("create users json decode err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := newUser.Validate(); err != nil {
		logger.Print(fmt.Sprintf("create users validation err=%s", err), ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dao.Create(newUser); err != nil {
		logger.Print(fmt.Sprintf("create users from dao err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	dao, err := elasticsearch.NewDao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userId := getParam("id", r); userId == "" {
		logger.Print(fmt.Sprintf("missing params of id, params=%s", userId), ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var updatedUser models.User
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

	if err := dao.Update(updatedUser); err != nil {
		logger.Print(fmt.Sprintf("update users dao err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	dao, err := elasticsearch.NewDao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := ""
	if userId = getParam("id", r); userId == "" {
		logger.Print(fmt.Sprintf("missing params of id, params=%s", userId), ERROR)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := dao.Delete(userId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
