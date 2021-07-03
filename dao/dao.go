package dao

import (
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/metildachee/userie/logger"
	"github.com/metildachee/userie/model"
)

type Dao struct {
	cli *elasticsearch.Client
}

func (dao *Dao) Init() error {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Print(fmt.Sprintf("failed to init es client %s", err), ERROR)
		return err
	}
	dao.cli = es

	res, err := dao.cli.Info()
	if err != nil {
		logger.Print(fmt.Sprintf("failed to get init info from es end point %s", err), ERROR)
		return err
	}

	logger.Print(fmt.Sprintf("init es client successfully %s", res), INFO)
	return nil
}

func (dao *Dao) GetUsers() (users []model.User, err error) {
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

func (dao *Dao) GetUser(userId int32) (user model.User, err error) {
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

func (dao *Dao) CreateUser(new model.User) (err error) {
	return
	// todo: upsert into es
}

func (dao *Dao) UpdateUser(id int32, new model.User) (err error) {
	// todo: upsert into es
	return
}

func (dao *Dao) DeleteUser(id int32) (err error) {
	// todo: remove from es
	return
}
