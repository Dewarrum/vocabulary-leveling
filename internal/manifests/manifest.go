package manifests

import (
	"dewarrum/vocabulary-leveling/internal/mpd"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type DbManifest struct {
	Id      uuid.UUID      `db:"id"`
	VideoId uuid.UUID      `db:"video_id"`
	Meta    types.JSONText `db:"meta"`
}

func NewDbManifest(videoId uuid.UUID, m *mpd.MPD) (*DbManifest, error) {
	metaJson, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &DbManifest{
		Id:      uuid.New(),
		VideoId: videoId,
		Meta:    types.JSONText(metaJson),
	}, nil
}

type DbManifestMeta struct {
	Profiles           string      `json:"profiles"`
	Type               string      `json:"type"`
	MaxSegmentDuration string      `json:"maxSegmentDuration"`
	MinBufferTime      string      `json:"minBufferTime"`
	Periods            []*DbPeriod `json:"periods"`
}

func newDbManifestMeta(m *mpd.MPD) *DbManifestMeta {
	return &DbManifestMeta{
		Profiles:           m.Profiles,
		Type:               m.Type,
		MaxSegmentDuration: m.MediaPresentationDuration,
		MinBufferTime:      m.MinBufferTime,
		Periods:            newPeriods(m.Periods),
	}
}

func newPeriods(periods []*mpd.Period) []*DbPeriod {
	var dbPeriods []*DbPeriod
	for _, period := range periods {
		dbPeriods = append(dbPeriods, &DbPeriod{
			Id:             period.ID,
			Start:          period.Start,
			AdaptationSets: newAdaptationSets(period.AdaptationSets),
		})
	}

	return dbPeriods
}

func newAdaptationSets(adaptationSets []*mpd.AdaptationSet) []*DbAdaptationSet {
	var dbAdaptationSets []*DbAdaptationSet
	for _, adaptationSet := range adaptationSets {
		dbAdaptationSets = append(dbAdaptationSets, &DbAdaptationSet{
			Id:                 adaptationSet.Id,
			ContentType:        adaptationSet.ContentType,
			StartWithSap:       adaptationSet.StartWithSAP,
			SegmentAlignment:   adaptationSet.SegmentAlignment,
			BitstreamSwitching: adaptationSet.BitstreamSwitching,
			Lang:               adaptationSet.Lang,
			MaxWidth:           adaptationSet.MaxWidth,
			MaxHeight:          adaptationSet.MaxHeight,
			Par:                adaptationSet.Par,
			Representations:    newRepresentations(adaptationSet.Representations),
		})
	}

	return dbAdaptationSets
}

func newRepresentations(representations []*mpd.Representation) []*DbRepresentation {
	var dbRepresentations []*DbRepresentation
	for _, representation := range representations {
		dbRepresentation := &DbRepresentation{
			Id:                representation.ID,
			MimeType:          representation.MimeType,
			Codecs:            representation.Codecs,
			Bandwidth:         representation.Bandwidth,
			AudioSamplingRate: representation.AudioSamplingRate,
			Width:             &representation.Width,
			Height:            &representation.Height,
			Sar:               &representation.SAR,
			SegmentTemplate:   newSegmentTemplate(representation.SegmentTemplate),
		}

		if representation.AudioChannelConfiguration != nil {
			dbRepresentation.AudioChannelConfiguration = newAudioChannelConfiguration(representation.AudioChannelConfiguration)
		}
		dbRepresentations = append(dbRepresentations, dbRepresentation)

	}

	return dbRepresentations
}

func newAudioChannelConfiguration(audioChannelConfiguration *mpd.AudioChannelConfiguration) *DbAudioChannelConfiguration {
	return &DbAudioChannelConfiguration{
		SchemeUriId: audioChannelConfiguration.SchemeIdUri,
		Value:       audioChannelConfiguration.Value,
	}
}

func newSegmentTemplate(segmentTemplate *mpd.SegmentTemplate) *DbSegmentTemplate {
	return &DbSegmentTemplate{
		Timescale: segmentTemplate.Timescale,
	}
}
