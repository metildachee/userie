package dao

import (
	"fmt"
	"time"

	"github.com/metildachee/userie/model"
)

func GetUsers() []model.User {
	// todo: access es and get all information of users
	users := make([]model.User, 0)

	for i := 0; i < 10; i++ {
		user := model.User{
			ID:          int32(i),
			Name:        fmt.Sprintf("metchee %d", i),
			DOB:         int32(time.Now().UnixNano()),
			Address:     "Kent Ridge",
			Description: "default user",
			Ctime:       int32(time.Now().UnixNano()),
		}
		users = append(users, user)
	}

	return users
}
