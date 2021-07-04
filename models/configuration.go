package models

import (
	"os"

	"github.com/google/logger"
)

type Tracer struct {
	ServiceName string `yaml:"service_name"`
}

type Configuration struct {
	ElasticEndpoint string `yaml:"elastic_endpoint"`
	ClusterName     string `yaml:"cluster_name"`
	ServerPort      string `yaml:"server_port"`
	Tracer          `yaml:"tracer"`
}

func (config *Configuration) Validate() bool {
	if config.ServiceName == "" {
		logger.Info("no service_name found but it is ok")
	}
	if config.ServerPort == "" {
		logger.Error("err config file missing server port")
		return false
	}
	if config.ClusterName == "" {
		logger.Error("err config file missing cluster name")
	}
	if config.ElasticEndpoint == "" {
		logger.Error("err config file missing elastic end point")
	}
	return true
}

func (config *Configuration) GetElasticEndpointEnvName() string {
	return "elastic_endpoint"
}

func (config *Configuration) GetClusterNameEnvName() string {
	return "cluster_name"
}

func (config *Configuration) GetServerEnvName() string {
	return "server_port"
}

func (config *Configuration) GetServiceEnvName() string {
	return "service_name"
}

func (config *Configuration) GetElasticEndpoint() string {
	if env := os.Getenv(config.GetElasticEndpointEnvName()); env != "" {
		return env
	}
	logger.Info("cannot get elastic end point from env, using default")
	return "http://127.0.0.1:9200"
}

func (config *Configuration) GetClusterName() string {
	if env := os.Getenv(config.GetClusterNameEnvName()); env != "" {
		return env
	}
	logger.Info("cannot get cluster name from env, using default")
	return "usersg0"
}

func (config *Configuration) GetServerEndpoint() string {
	if env := os.Getenv(config.GetServerEnvName()); env != "" {
		return env
	}
	logger.Info("cannot get port from env, using default")
	return ":8080"
}

func (config *Configuration) GetServiceName() string {
	if env := os.Getenv(config.GetServiceEnvName()); env != "" {
		return env
	}
	logger.Info("cannot get service name from env, using default")
	return "userie (default)"
}
