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
	users, err := dao.GetUsers()
	if err != nil {
		// todo: return err code in response
		log.Print(fmt.Sprintf("get users from dao err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Print(fmt.Sprintf("get users json encoder err=%s", err), ERROR)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
