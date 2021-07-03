package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	query  map[string]interface{}
	match  map[string]interface{}
	result map[string]interface{}
}

func NewDao() (*Dao, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Print(fmt.Sprintf("failed to init es client %s", err), ERROR)
		return nil, err
	}
	dao := &Dao{}
	m := make(map[string]interface{})
	dao.cli = es
	dao.cluster = "usersg0"
	dao.builder.query = m
	dao.builder.match = m
	dao.builder.result = m
	logger.Print(fmt.Sprintf("init es client successfull"), INFO)
	return dao, nil
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
	return
}

func (dao *Dao) GetUser(id string) (user model.User, err error) {
	logger.Print("getting user", INFO)
	query, err := dao.AddMatch("_id", id).Build()
	fmt.Println("query", query)
	if err != nil {
		return
	}

	res, err := dao.Search(query)

	if err != nil {
		logger.Print(fmt.Sprintf("search err=%s", err), ERROR)
		return
	}
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			logger.Print(fmt.Sprintf("Error parsing the response body: %s", err), ERROR)
		} else {
			// Print the response status and error information.
			logger.Print(fmt.Sprintf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"]), ERROR)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&dao.builder.result); err != nil {
		logger.Print(fmt.Sprintf("json unmarshal err=%s", err), ERROR)
	}
	logger.Print("got results from es", INFO)
	for _, hit := range dao.builder.result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		logger.Print(fmt.Sprintf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"]), INFO)
		if hit.(map[string]interface{})["_id"] == "" {
			continue
		}
		user.ID = hit.(map[string]interface{})["_id"].(string)
		user.Ctime = (hit.(map[string]interface{})["_source"].(map[string]interface{}))["ctime"].(float64)
		user.Name = (hit.(map[string]interface{})["_source"].(map[string]interface{}))["name"].(string)
		user.DOB = (hit.(map[string]interface{})["_source"].(map[string]interface{}))["dob"].(float64)
		user.Description = (hit.(map[string]interface{})["_source"].(map[string]interface{}))["description"].(string)
		user.Description = (hit.(map[string]interface{})["_source"].(map[string]interface{}))["address"].(string)
		break
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
	logger.Print(fmt.Sprintf("created doc successfully"), INFO)
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
	logger.Print(fmt.Sprintf("searching in %s", dao.cluster), INFO)
	return dao.cli.Search(
		dao.cli.Search.WithContext(context.Background()),
		dao.cli.Search.WithIndex(dao.cluster),
		dao.cli.Search.WithBody(buf),
		dao.cli.Search.WithTrackTotalHits(true),
		dao.cli.Search.WithPretty())
}

func (dao *Dao) Build() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	query := make(map[string]interface{})
	query["query"] = dao.builder.match
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Print(fmt.Sprintf("build query failed %s", err), ERROR)
		return nil, err
	}
	return &buf, nil
}

func (dao *Dao) AddMatch(field string, val interface{}) *Dao {
	logger.Print("add match to builder", INFO)
	q := make(map[string]interface{})
	q[field] = val
	dao.builder.match["match"] = q
	return dao
}
