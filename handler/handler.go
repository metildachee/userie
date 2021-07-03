package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/metildachee/userie/dao"
	"github.com/metildachee/userie/log"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users := dao.GetUsers()
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Print(fmt.Sprintf("get users unexpected err=%s", err), INFO)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

}
