package subtitles

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

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

type FtsSubtitleCue struct {
	Id       uuid.UUID `json:"cue_id"`
	VideoId  uuid.UUID `json:"video_id"`
	Sequence int       `json:"sequence"`
	Text     string    `json:"text"`
}

func NewFtsSubtitleCue(id, videoId uuid.UUID, sequence int, text string) *FtsSubtitleCue {
	return &FtsSubtitleCue{
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

func (f *FullTextSearch) Insert(ftsSubtitleCue *FtsSubtitleCue, context context.Context) error {
	documentId := fmt.Sprintf("%s/%d", ftsSubtitleCue.VideoId, ftsSubtitleCue.Sequence)
	_, err := f.elasticsearchClient.Index(indexName).
		Id(base64.URLEncoding.EncodeToString([]byte(documentId))).
		Request(ftsSubtitleCue).
		Do(context)

	if err == nil {
		return nil
	}

	elasticsearchErr := err.(*types.ElasticsearchError)
	f.logger.Error().Str("videoId", ftsSubtitleCue.VideoId.String()).Int32("sequence", int32(ftsSubtitleCue.Sequence)).Err(elasticsearchErr).Msg("Failed to insert subtitle cue")

	return errors.Join(err, ErrFailedToInsert)
}

func (f *FullTextSearch) Search(queryText string, context context.Context) ([]*FtsSubtitleCue, error) {
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

	var subtitleCues []*FtsSubtitleCue
	for _, hit := range response.Hits.Hits {
		subtitleCue := FtsSubtitleCue{}
		err := json.Unmarshal(hit.Source_, &subtitleCue)
		if err != nil {
			return nil, errors.Join(err, ErrFailedToDeserialize)
		}
		subtitleCues = append(subtitleCues, &subtitleCue)
	}

	return subtitleCues, nil
}

func initializeIndex(elasticsearchClient *elasticsearch.TypedClient, context context.Context) error {
	_, err := elasticsearchClient.Indices.
		Create(indexName).
		Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"cue_id":   types.KeywordProperty{},
					"video_id": types.KeywordProperty{},
					"sequence": types.KeywordProperty{},
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
