package models

import "os"

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
	if config.ServerPort == "" || config.ClusterName == "" || config.ElasticEndpoint == "" {
		return false
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
	return "http://127.0.0.1:9200"
}

func (config *Configuration) GetClusterName() string {
	if env := os.Getenv(config.GetClusterNameEnvName()); env != "" {
		return env
	}
	return "usersg0"
}

func (config *Configuration) GetServerEndpoint() string {
	if env := os.Getenv(config.GetServerEnvName()); env != "" {
		return env
	}
	return ":8080"
}

func (config *Configuration) GetServiceName() string {
	if env := os.Getenv(config.GetServiceEnvName()); env != "" {
		return env
	}
	return "userie (default)"
}