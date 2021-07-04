package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/google/logger"
	"github.com/metildachee/userie/models"
	elasticv7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type UserImplDao struct {
	cli     *elasticv7.Client
	cluster string
	safe    SafeCounter
}

func (dao *UserImplDao) GetAll(ctx context.Context, limit int) (users []models.User, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "es get all")
	defer span.Finish()

	if !dao.CheckInit(ctx) {
		return users, errors.New("es client not init")
	}
	query := elasticv7.NewBoolQuery().
		Must(elasticv7.NewExistsQuery("id"))
	src, err := query.Source()
	span.LogFields(log.String("es query", fmt.Sprintf("%v", src)))

	if err != nil {
		ext.LogError(span, err)
	}
	span.LogFields(log.String("limit", string(rune(limit))))
	searchResult, err := dao.cli.Search().
		Index(dao.cluster).
		Query(query).
		From(0).
		Size(limit).
		Do(ctx)
	if err != nil {
		ext.LogError(span, err)
		return
	}
	for _, item := range searchResult.Each(reflect.TypeOf(models.User{})) {
		if u, ok := item.(models.User); ok {
			users = append(users, u)
		}
	}
	span.LogFields(log.String("users", fmt.Sprintf("%v", users)))
	return
}

func (dao *UserImplDao) GetById(ctx context.Context, id string) (user models.User, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "es by id")
	defer span.Finish()

	if !dao.CheckInit(ctx) {
		return user, errors.New("es client not init")
	}
	query := elasticv7.NewTermQuery("id", id)
	src, err := query.Source()
	if err != nil {
		ext.LogError(span, err)
	}
	span.LogFields(log.String("es query", fmt.Sprintf("%v", src)))

	searchResult, err := dao.cli.Search().
		Index(dao.cluster).
		Query(query).
		Do(ctx)
	if err != nil {
		ext.LogError(span, err)
		return
	}
	for _, item := range searchResult.Each(reflect.TypeOf(user)) {
		if u, ok := item.(models.User); ok {
			user = u
			span.LogFields(log.String("user", fmt.Sprintf("%v", user)))
			return
		}
	}
	return user, errors.New("nil hit")
}

func (dao *UserImplDao) Create(ctx context.Context, new models.User, wg ...*sync.WaitGroup) (id string, err error) {
	if len(wg) > 0 {
		defer wg[0].Done()
	}
	if !dao.CheckInit(ctx) {
		return id, errors.New("es client not init")
	}
	if id, err = dao.create(ctx, new); err != nil {
		return
	}
	return
}

func (dao *UserImplDao) create(ctx context.Context, new models.User, wg ...*sync.WaitGroup) (id string, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "es create item")
	defer span.Finish()

	if len(wg) > 0 {
		defer wg[0].Done()
	}

	new.ID = dao.safe.GetCount()
	doc, err := json.Marshal(new)
	if err != nil {
		ext.LogError(span, err)
		fmt.Println("json marshal err", err)
		return
	}

	put1, err := dao.cli.Index().
		Index(dao.cluster).
		Id(new.ID).
		BodyJson(string(doc)).
		Do(ctx)
	if err != nil {
		ext.LogError(span, err)
		fmt.Println("index document failed", err, "cluster name", dao.cluster)
		return
	}

	_, err = dao.cli.Flush().
		Index(dao.cluster).
		Do(ctx)
	if err != nil {
		ext.LogError(span, err)
		fmt.Println("flushing index failed", err)
		return
	}

	id = put1.Id
	span.LogFields(
		log.String("user doc", put1.Id),
		log.String("user index", put1.Index))
	return
}

func (dao *UserImplDao) BatchCreate(ctx context.Context, new []models.User) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "es batch item")
	defer span.Finish()

	if !dao.CheckInit(ctx) {
		return errors.New("es client not init")
	}
	var wg sync.WaitGroup

	for _, item := range new {
		wg.Add(1)
		go dao.create(ctx, item, &wg)
	}
	wg.Wait()

	span.LogKV("batch index done")
	return
}

func (dao *UserImplDao) Update(ctx context.Context, updated models.User) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "es update item")
	defer span.Finish()

	if !dao.CheckInit(ctx) {
		return errors.New("es client not init")
	}

	doc, err := json.Marshal(updated)
	if err != nil {
		logger.Error(err)
		ext.LogError(span, err)
		return
	}
	update, err := dao.cli.Index().
		Index(dao.cluster).
		Id(updated.ID).
		BodyJson(string(doc)).
		Do(ctx)
	if err != nil {
		ext.LogError(span, err)
		logger.Error(err)
		return
	}
	span.LogFields(
		log.String("user doc", update.Id),
		log.String("user index", strconv.FormatInt(update.Version, 10)))
	return
}

func (dao *UserImplDao) UpdateUserName(ctx context.Context, id, newName string) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "es update name of item")
	defer span.Finish()

	span.LogFields(
		log.String("id", id),
		log.String("new name", newName))
	if !dao.CheckInit(ctx) {
		return errors.New("es client not init")
	}
	update, err := dao.cli.Update().
		Index(dao.cluster).
		Id(id).
		Doc(map[string]interface{}{"name": newName}).
		Do(ctx)
	if err != nil {
		ext.LogError(span, err)
		return
	}
	span.LogFields(
		log.String("user doc", update.Id),
		log.String("user index", strconv.FormatInt(update.Version, 10)))
	return
}

func (dao *UserImplDao) Delete(ctx context.Context, id string) (err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "es delete item")
	defer span.Finish()
	span.LogFields(log.String("doc id", id))

	if !dao.CheckInit(ctx) {
		return errors.New("es client not init")
	}
	_, err = dao.cli.Delete().
		Index(dao.cluster).
		Id(id).Refresh("true").
		Do(ctx)
	if err != nil {
		ext.LogError(span, err)
		return
	}
	return
}
