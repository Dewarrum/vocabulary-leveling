package app

import (
	"errors"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	ErrElasticsearchUrlIsRequired      = errors.New("ELASTICSEARCH_URL is required")
	ErrElasticsearchUsernameIsRequired = errors.New("ELASTICSEARCH_USERNAME is required")
	ErrElasticsearchPasswordIsRequired = errors.New("ELASTICSEARCH_PASSWORD is required")
)

func createElasticSearchClient() (*elasticsearch.TypedClient, error) {
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		return nil, ErrElasticsearchUrlIsRequired
	}

	username := os.Getenv("ELASTICSEARCH_USERNAME")
	if username == "" {
		return nil, ErrElasticsearchUsernameIsRequired
	}

	password := os.Getenv("ELASTICSEARCH_PASSWORD")
	if password == "" {
		return nil, ErrElasticsearchPasswordIsRequired
	}
	elasticsearchClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{url},
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}

	return elasticsearchClient, nil
}
