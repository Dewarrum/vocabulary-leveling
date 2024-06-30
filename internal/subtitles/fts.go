package subtitles

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	indexName = "subtitles"
)

var (
	ErrFailedToCreateIndex = errors.New("failed to create index")
	ErrFailedToInsert      = errors.New("failed to insert")
	ErrFailedToSearch      = errors.New("failed to search")
	ErrFailedToDeserialize = errors.New("failed to deserialize")
	analyzer               = "nori"
)

type FtsSubtitle struct {
	Id       string    `json:"id"`
	VideoId  uuid.UUID `json:"video_id"`
	Sequence int       `json:"sequence"`
	Text     string    `json:"text"`
}

func NewFtsSubtitle(id string, videoId uuid.UUID, sequence int, text string) *FtsSubtitle {
	return &FtsSubtitle{
		Id:       id,
		VideoId:  videoId,
		Sequence: sequence,
		Text:     text,
	}
}

type FullTextSearch struct {
	elasticsearchClient *elasticsearch.TypedClient
	logger              zerolog.Logger
}

func NewFullTextSearch(elasticsearchClient *elasticsearch.TypedClient, context context.Context) (*FullTextSearch, error) {
	err := initializeIndex(elasticsearchClient, context)
	if err != nil {
		return nil, errors.Join(err, ErrFailedToCreateIndex)
	}

	return &FullTextSearch{
		elasticsearchClient: elasticsearchClient,
	}, nil
}

func (f *FullTextSearch) Insert(subtitle *FtsSubtitle, context context.Context) error {
	_, err := f.elasticsearchClient.Index(indexName).
		Request(subtitle).
		Do(context)

	if err == nil {
		return nil
	}

	elasticsearchErr := err.(*types.ElasticsearchError)
	f.logger.Error().Str("videoId", subtitle.VideoId.String()).Int32("sequence", int32(subtitle.Sequence)).Err(elasticsearchErr).Msg("Failed to insert subtitle")

	return errors.Join(err, ErrFailedToInsert)
}

func (f *FullTextSearch) Search(queryText string, context context.Context) ([]*FtsSubtitle, error) {
	query := types.Query{
		Match: map[string]types.MatchQuery{
			"text": {
				Query: queryText,
			},
		},
	}

	response, err := f.elasticsearchClient.Search().
		Index(indexName).
		Query(&query).
		Do(context)

	if err != nil {
		return nil, errors.Join(err, ErrFailedToSearch)
	}

	var subtitles []*FtsSubtitle
	for _, hit := range response.Hits.Hits {
		subtitle := FtsSubtitle{}
		err := json.Unmarshal(hit.Source_, &subtitle)
		if err != nil {
			return nil, errors.Join(err, ErrFailedToDeserialize)
		}
		subtitles = append(subtitles, &subtitle)
	}

	return subtitles, nil
}

func initializeIndex(elasticsearchClient *elasticsearch.TypedClient, context context.Context) error {
	_, err := elasticsearchClient.Indices.
		Create(indexName).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"subtitle_id": types.KeywordProperty{},
					"video_id":    types.KeywordProperty{},
					"sequence":    types.KeywordProperty{},
					"text": types.TextProperty{
						Analyzer: &analyzer,
					},
				},
			},
		}).
		Do(context)

	if err == nil {
		return nil
	}

	elasticsearchErr := err.(*types.ElasticsearchError)
	if elasticsearchErr.ErrorCause.Type == "resource_already_exists_exception" {
		return nil
	}

	return err
}
