package dao

import (
	"fmt"
	"time"

	"github.com/metildachee/userie/model"
)

func GetUsers() (users []model.User, err error) {
	// todo: access es and get all information of users

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
	return
}

func GetUser(userId int32) (user model.User, err error) {
	// todo: access es and get all information of users
	return model.User{
		ID:          1,
		Name:        "metchee",
		DOB:         1625276913,
		Address:     "Kent Ridge",
		Description: "default user description",
		Ctime:       1625276913,
	}, nil
}

func CreateUser(new model.User) (err error) {
	return
	// todo: upsert into es and
}

func UpdateUser(id int32, new model.User) (err error) {
	// todo: upsert into es
	return
}

func DeleteUser(id int32) (err error) {
	return
}
