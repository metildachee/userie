package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/metildachee/userie/logger"
	"github.com/metildachee/userie/model"
)

type Dao struct {
	cli     *elasticsearch.Client
	cluster string
}

func NewDao() *Dao {
	return &Dao{}
}

func (dao *Dao) Init() error {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Print(fmt.Sprintf("failed to init es client %s", err), ERROR)
		return err
	}
	dao.cli = es
	dao.cluster = "usersg0"
	logger.Print(fmt.Sprintf("init es client successfull"), INFO)
	return nil
}

func (dao *Dao) CheckInit() bool {
	if dao.cli == nil {
		return false
	}
	_, err := dao.cli.Info()
	if err != nil {
		logger.Print(fmt.Sprintf("es client is not inited %s", err), INFO)
		return false
	}
	logger.Print(fmt.Sprintf("es client is inited"), INFO)
	return true
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
	if !dao.CheckInit() {
		logger.Print(fmt.Sprintf("es client is not inited %s", err), ERROR)
		return errors.New("es client not initialised")
	}

	req, err := dao.GetIndexRequest(new)
	if err != nil {
		logger.Print(fmt.Sprintf("%s", err), ERROR)
		return
	}

	res, err := req.Do(context.Background(), dao.cli)
	if err != nil {
		logger.Print(fmt.Sprintf("%s", err), ERROR)
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		logger.Print(fmt.Sprintf("[%s] error indexing doc", res.Status()), ERROR)
		return errors.New("error while indexing document")
	}
	return
}

func (dao *Dao) GetIndexRequest(doc model.User) (*esapi.IndexRequest, error) {
	userDoc, err := json.Marshal(doc)
	if err != nil {
		logger.Print(fmt.Sprintf("json marshal failed=%s", err), ERROR)
		return nil, err
	}

	return &esapi.IndexRequest{
		Index:   dao.cluster,
		Body:    strings.NewReader(string(userDoc)),
		Refresh: "true",
	}, nil
}

func (dao *Dao) UpdateUser(id int32, new model.User) (err error) {
	// todo: upsert into es
	return
}

func (dao *Dao) DeleteUser(id int32) (err error) {
	// todo: remove from es
	return
}
