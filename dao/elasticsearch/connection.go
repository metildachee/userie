package elasticsearch

import (
	"context"
	"fmt"
	"os"

	"github.com/google/logger"
	"github.com/metildachee/userie/models"
	elasticv7 "github.com/olivere/elastic/v7"
)

func NewDao() (*UserImplDao, error) {
	config := models.Configuration{}
	es, err := elasticv7.NewSimpleClient(elasticv7.SetURL(config.GetElasticEndpoint()))
	if err != nil {
		logger.Info(fmt.Sprintf("failed to init es client %s", err), ERROR)
		return nil, err
	}

	dao := &UserImplDao{}
	dao.cli = es
	dao.cluster = os.Getenv(config.GetClusterName())
	dao.ctx = context.Background()
	logger.Info(fmt.Sprintf("init es client successfull"), INFO)
	return dao, nil
}

func (dao *UserImplDao) CheckInit() bool {
	if dao.cli == nil {
		return false
	}
	exists, err := dao.cli.IndexExists(dao.cluster).Do(dao.ctx)
	if err != nil {
		logger.Info(fmt.Sprintf("error when getting index exists %s", err), ERROR)
		return false
	}
	if !exists {
		logger.Info(fmt.Sprintf("es index is not inited %s", err), ERROR)
		return false
	}
	logger.Info(fmt.Sprintf("es client is inited"), INFO)
	return true
}
