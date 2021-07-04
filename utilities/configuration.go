package utilities

import (
	"errors"
	"os"

	"github.com/google/logger"
	"github.com/metildachee/userie/models"
	"gopkg.in/yaml.v3"
)

func SetConfig(configPath string) (config models.Configuration, err error) {
	file, err := os.Open(configPath)
	if err != nil {
		logger.Error("err opening up config path", err)
		return
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		logger.Error("error while decoding configuration file", configPath)
		return
	}

	if !config.Validate() {
		logger.Error("config is invalid")
		return config, errors.New("invalid configuration file")
	}

	os.Setenv(config.GetElasticEndpointEnvName(), config.ElasticEndpoint)
	os.Setenv(config.GetClusterNameEnvName(), config.ClusterName)
	os.Setenv(config.GetServerEnvName(), config.ServerPort)

	logger.Info("set config successfully")
	return
}
