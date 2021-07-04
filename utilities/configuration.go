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

	if !valid(config) {
		return config, errors.New("invalid configuration file")
	}

	os.Setenv("ELASTIC_ENDPOINT", config.ElasticEndpoint)
	os.Setenv("CLUSTER_NAME", config.ClusterName)
	os.Setenv("SERVER_PORT", config.ServerPort)
	return
}

func valid(config models.Configuration) bool {
	if config.ServerPort == "" || config.ClusterName == "" || config.ElasticEndpoint == "" {
		return false
	}
	return true
}
