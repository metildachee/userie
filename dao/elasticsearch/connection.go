package elasticsearch

import (
	"context"
	"errors"
	"os"

	"github.com/google/logger"
	"github.com/metildachee/userie/models"
	elasticv7 "github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func NewDao(ctx context.Context) (*UserImplDao, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "get new dao")
	defer span.Finish()

	config := models.Configuration{}
	es, err := elasticv7.NewSimpleClient(elasticv7.SetURL(config.GetElasticEndpoint()))
	if err != nil {
		ext.LogError(span, err)
		return nil, err
	}

	dao := &UserImplDao{}
	dao.cli = es
	dao.cluster = os.Getenv(config.GetClusterName())
	span.LogKV("es client init successfully")
	return dao, nil
}

func (dao *UserImplDao) CheckInit(ctx context.Context) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "check dao status")
	defer span.Finish()

	if dao.cli == nil {
		ext.LogError(span, errors.New("es client does not exist"))
		logger.Fatalf("es client does not exist, exiting")
		return false
	}
	exists, err := dao.cli.IndexExists(dao.cluster).Do(ctx)
	if err != nil {
		ext.LogError(span, err)
		return false
	}
	if !exists {
		ext.LogError(span, errors.New("index does not exists"))
		logger.Fatalf("index does not exists, exiting")
		return false
	}
	span.LogKV("es client is ok")
	logger.Info("es client is ok")
	return true
}
