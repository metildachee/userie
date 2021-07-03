package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/metildachee/userie/logger"
	"github.com/metildachee/userie/model"
	elasticv7 "github.com/olivere/elastic/v7"
)

type Dao struct {
	cli     *elasticv7.Client
	cluster string
}

func NewDao() (*Dao, error) {
	es, err := elasticv7.NewSimpleClient(elasticv7.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		logger.Print(fmt.Sprintf("failed to init es client %s", err), ERROR)
		return nil, err
	}
	dao := &Dao{}
	dao.cli = es
	dao.cluster = "usersg0"
	logger.Print(fmt.Sprintf("init es client successfull"), INFO)
	return dao, nil
}

func (dao *Dao) CheckInit() bool {
	ctx := context.Background()
	if dao.cli == nil {
		return false
	}
	exists, err := dao.cli.IndexExists(dao.cluster).Do(ctx)
	if err != nil {
		logger.Print(fmt.Sprintf("error when getting index exists %s", err), ERROR)
		return false
	}

	if !exists {
		logger.Print(fmt.Sprintf("es index is not inited %s", err), ERROR)
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
	ctx := context.Background()
	logger.Print("getting user", INFO)
	if !dao.CheckInit() {
		return user, errors.New("es client not init")
	}
	query := elasticv7.NewMatchQuery("_id", id)
	searchResult, err := dao.cli.Search().Index(dao.cluster).Query(query).Do(ctx)
	if err != nil {
		logger.Print(fmt.Sprintf("search err=%s", err), ERROR)
		return
	}
	for _, item := range searchResult.Each(reflect.TypeOf(user)) {
		if u, ok := item.(model.User); ok {
			user = u
			fmt.Printf("user details by %s: %s\n", u.ID, u.Name)
			return
		}
	}
	return
}

func (dao *Dao) CreateUser(new model.User) (err error) {
	if !dao.CheckInit() {
		return errors.New("es client not init")
	}

	ctx := context.Background()
	doc, err := json.Marshal(new)
	if err != nil {
		return err
	}
	fmt.Println(string(doc))
	put1, err := dao.cli.Index().Index(dao.cluster).Id("").BodyJson(string(doc)).Do(ctx)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	_, err = dao.cli.Flush().Index(dao.cluster).Do(ctx)
	if err != nil {
		panic(err)
	}

	logger.Print(fmt.Sprintf("Indexed user doc %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type), INFO)
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
