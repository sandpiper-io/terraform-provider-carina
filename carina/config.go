package carina

import (
	"github.com/getcarina/libcarina"
)

type Config struct {
	Username    string
	ApiKey string
}

func (c *Config) NewClient() (*libcarina.ClusterClient, error) {

	client, err := libcarina.NewClusterClient(libcarina.BetaEndpoint, c.Username, c.ApiKey)

	return client, err
}