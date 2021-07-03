package dao

import (
	"bytes"
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
	builder Query
}

type Query struct {
	query map[string]interface{}
	match map[string]interface{}
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
	query, err := dao.AddMatch("_id", userId).Build()
	if err != nil {
		return
	}

	res, err := dao.Search(query)
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			logger.Print(fmt.Sprintf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"]), ERROR)
		}
	}

	if err != nil {
		logger.Print(fmt.Sprintf("search %s", err), ERROR)
		return
	}
	return
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

func (dao *Dao) UpdateUser(id int32, new model.User) (err error) {
	// todo: upsert into es
	return
}

func (dao *Dao) DeleteUser(id int32) (err error) {
	// todo: remove from es
	return
}

///

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

func (dao *Dao) Search(buf *bytes.Buffer) (*esapi.Response, error) {
	return dao.cli.Search(
		dao.cli.Search.WithContext(context.Background()),
		dao.cli.Search.WithIndex(dao.cluster),
		dao.cli.Search.WithBody(buf),
		dao.cli.Search.WithTrackTotalHits(true),
		dao.cli.Search.WithPretty())
}

func (dao *Dao) Build() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	dao.builder.query["query"] = dao.builder.match
	if err := json.NewEncoder(&buf).Encode(dao.builder.query); err != nil {
		logger.Print(fmt.Sprintf("build query failed %s", err), ERROR)
		return nil, err
	}
	return &buf, nil
}

func (dao *Dao) AddMatch(field string, val interface{}) *Dao {
	dao.builder.match[field] = val
	return dao
}
