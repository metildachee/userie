package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/metildachee/userie/logger"
	"github.com/metildachee/userie/model"
	elasticv7 "github.com/olivere/elastic/v7"
)

type Dao struct {
	cli     *elasticv7.Client
	cluster string
	SafeCounter
	ctx context.Context
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
	dao.ctx = context.Background()
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

func (dao *Dao) GetUsers(limit int) (users []model.User, err error) {
	if !dao.CheckInit() {
		return users, errors.New("es client not init")
	}
	query := elasticv7.NewBoolQuery().
		Must(elasticv7.NewExistsQuery("id")) // see: https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-exists-query.html

	searchResult, err := dao.cli.Search().
		Index(dao.cluster).
		Query(query).
		From(0).
		Size(limit).
		Do(dao.ctx)
	if err != nil {
		logger.Print(fmt.Sprintf("search err=%s", err), ERROR)
		return
	}
	for _, item := range searchResult.Each(reflect.TypeOf(model.User{})) {
		if u, ok := item.(model.User); ok {
			fmt.Printf("user details by %s: %s\n", u.ID, u.Name)
			users = append(users, u)
		}
	}
	return
}

func (dao *Dao) GetUser(id string) (user model.User, err error) {
	if !dao.CheckInit() {
		return user, errors.New("es client not init")
	}
	query := elasticv7.NewTermQuery("id", id) // see: https://www.elastic.co/guide/en/elasticsearch/reference/6.8/query-dsl-term-query.html
	searchResult, err := dao.cli.Search().
		Index(dao.cluster).
		Query(query).
		Do(dao.ctx)
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
	return user, errors.New("nil hit")
}

func (dao *Dao) CreateUser(new model.User, wg ...*sync.WaitGroup) (err error) {
	if len(wg) > 0 {
		defer wg[0].Done()
	}
	if !dao.CheckInit() {
		return errors.New("es client not init")
	}
	if err = dao.createUser(new); err != nil {
		return
	}
	return
}

/* creates user straight away without checking if there are clients */
func (dao *Dao) createUser(new model.User, wg ...*sync.WaitGroup) (err error) {
	if len(wg) > 0 {
		defer wg[0].Done()
	}

	new.ID = dao.GetCount()
	doc, err := json.Marshal(new)
	if err != nil {
		return err
	}
	put1, err := dao.cli.Index().Index(dao.cluster).Id(new.ID).BodyJson(string(doc)).Do(dao.ctx)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	_, err = dao.cli.Flush().Index(dao.cluster).Do(dao.ctx)
	if err != nil {
		panic(err)
	}

	logger.Print(fmt.Sprintf("Indexed user doc %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type), INFO)
	return
}

func (dao *Dao) BatchCreateUsers(new []model.User) (err error) {
	if !dao.CheckInit() {
		return errors.New("es client not init")
	}
	var wg sync.WaitGroup

	for _, item := range new {
		wg.Add(1)
		go dao.createUser(item, &wg)
	}
	wg.Wait()

	logger.Print("batch index done", INFO)
	return
}

func (dao *Dao) UpdateUser(updated model.User) (err error) {
	if !dao.CheckInit() {
		return errors.New("es client not init")
	}

	doc, err := json.Marshal(updated)
	if err != nil {
		logger.Print(fmt.Sprintf("err when marshalling json %s", err), ERROR)
		return
	}
	update, err := dao.cli.Index().
		Index(dao.cluster).
		Id(updated.ID).
		BodyJson(string(doc)).
		Do(dao.ctx)
	if err != nil {
		logger.Print(fmt.Sprintf("es error %s", err), ERROR)
		return
	}
	logger.Print(fmt.Sprintf("New version of user %q is now %d", update.Id, update.Version), INFO)
	return
}

func (dao *Dao) UpdateUserName(id, newName string) (err error) {
	if !dao.CheckInit() {
		return errors.New("es client not init")
	}
	update, err := dao.cli.Update().Index(dao.cluster).Id(id).Doc(map[string]interface{}{"name": newName}).Do(dao.ctx)
	if err != nil {
		logger.Print(fmt.Sprintf("es error %s", err), ERROR)
		return
	}
	logger.Print(fmt.Sprintf("New version of user %q is now %d", update.Id, update.Version), INFO)
	return
}

func (dao *Dao) DeleteUser(id string) (err error) {
	if !dao.CheckInit() {
		return errors.New("es client not init")
	}
	_, err = dao.cli.Delete().
		Index(dao.cluster).
		Id(id).Refresh("true").
		Do(dao.ctx)
	if err != nil {
		return
	}
	logger.Print("deleted successfully", INFO)
	return
}
