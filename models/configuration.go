package models

type Configuration struct {
	ElasticEndpoint string `yaml:"elastic_endpoint", envconfig:"ELASTIC_ENDPOINT"`
	ClusterName     string `yaml:"cluster_name", envconfig:"CLUSTER_NAME"`
	ServerPort      string `yaml:"server_port", envconfig:"SERVER_PORT"`
}
