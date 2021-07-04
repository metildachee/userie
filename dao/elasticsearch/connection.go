package elasticsearch

import (
	"context"
	"fmt"
	"os"

	"github.com/metildachee/userie/logger"
	elasticv7 "github.com/olivere/elastic/v7"
)

func NewDao() (*UserImplDao, error) {
	es, err := elasticv7.NewSimpleClient(elasticv7.SetURL(os.Getenv("ELASTIC_ENDPOINT")))
	if err != nil {
		logger.Print(fmt.Sprintf("failed to init es client %s", err), ERROR)
		return nil, err
	}
	dao := &UserImplDao{}
	dao.cli = es
	dao.cluster = os.Getenv("CLUSTER_NAME")
	dao.ctx = context.Background()
	logger.Print(fmt.Sprintf("init es client successfull"), INFO)
	return dao, nil
}

func (dao *UserImplDao) CheckInit() bool {
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
