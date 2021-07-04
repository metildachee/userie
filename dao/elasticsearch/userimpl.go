package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/metildachee/userie/logger"
	"github.com/metildachee/userie/models"
	elasticv7 "github.com/olivere/elastic/v7"
)

type UserImplDao struct {
	cli     *elasticv7.Client
	cluster string
	safe    SafeCounter
	ctx     context.Context
}

func (dao *UserImplDao) GetAll(limit int) (users []models.User, err error) {
	if !dao.CheckInit() {
		return users, errors.New("es client not init")
	}
	query := elasticv7.NewBoolQuery().
		Must(elasticv7.NewExistsQuery("id"))

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
	for _, item := range searchResult.Each(reflect.TypeOf(models.User{})) {
		if u, ok := item.(models.User); ok {
			fmt.Printf("user details by %s: %s\n", u.ID, u.Name)
			users = append(users, u)
		}
	}
	return
}

func (dao *UserImplDao) GetById(id string) (user models.User, err error) {
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
		if u, ok := item.(models.User); ok {
			user = u
			fmt.Printf("user details by %s: %s\n", u.ID, u.Name)
			return
		}
	}
	return user, errors.New("nil hit")
}

func (dao *UserImplDao) Create(new models.User, wg ...*sync.WaitGroup) (id string, err error) {
	if len(wg) > 0 {
		defer wg[0].Done()
	}
	if !dao.CheckInit() {
		return id, errors.New("es client not init")
	}
	if id, err = dao.create(new); err != nil {
		return
	}
	return
}

func (dao *UserImplDao) create(new models.User, wg ...*sync.WaitGroup) (id string, err error) {
	if len(wg) > 0 {
		defer wg[0].Done()
	}

	new.ID = dao.safe.GetCount()
	doc, err := json.Marshal(new)
	if err != nil {
		return
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

	id = put1.Id
	logger.Print(fmt.Sprintf("Indexed user doc %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type), INFO)
	return
}

func (dao *UserImplDao) BatchCreate(new []models.User) (err error) {
	if !dao.CheckInit() {
		return errors.New("es client not init")
	}
	var wg sync.WaitGroup

	for _, item := range new {
		wg.Add(1)
		go dao.create(item, &wg)
	}
	wg.Wait()

	logger.Print("batch index done", INFO)
	return
}

func (dao *UserImplDao) Update(updated models.User) (err error) {
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

func (dao *UserImplDao) UpdateUserName(id, newName string) (err error) {
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

func (dao *UserImplDao) Delete(id string) (err error) {
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
