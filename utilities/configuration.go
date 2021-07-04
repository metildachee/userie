package utilities

import (
	"errors"
	"os"

	"github.com/metildachee/userie/models"
	"gopkg.in/yaml.v3"
)

func SetConfig(configPath string) (config models.Configuration, err error) {
	file, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return
	}

	if !config.Validate(){
		return config, errors.New("invalid configuration file")
	}

	os.Setenv(config.GetElasticEndpointEnvName(), config.ElasticEndpoint)
	os.Setenv(config.GetClusterNameEnvName(), config.ClusterName)
	os.Setenv(config.GetServerEnvName(), config.ServerPort)
	return
}

