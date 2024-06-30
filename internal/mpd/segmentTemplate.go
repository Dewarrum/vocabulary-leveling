package mpd

import (
	"errors"
	"strconv"
)

type SegmentTemplate struct {
	Timescale       string           `xml:"timescale,attr" json:"timescale,omitempty"`
	Initialization  *Initialization  `xml:"Initialization" json:"initialization,omitempty"`
	Media           string           `xml:"media,attr" json:"media,omitempty"`
	Duration        string           `xml:"duration,attr" json:"duration,omitempty"`
	StartNumber     string           `xml:"startNumber,attr" json:"startNumber,omitempty"`
	Times           string           `xml:"times,attr" json:"times,omitempty"`
	Presentation    string           `xml:"presentation,attr" json:"presentation,omitempty"`
	Bandwidth       string           `xml:"bandwidth,attr" json:"bandwidth,omitempty"`
	ProgramDateTime string           `xml:"programDateTime,attr" json:"programDateTime,omitempty"`
	SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline,omitempty" json:"segmentTimeline,omitempty"`
}

type SegmentTemplateEntryInfo struct {
	TimestampMs int64
	DurationMs  int64
}

func (s *SegmentTemplate) GetSegmentInfos() ([]*SegmentTemplateEntryInfo, error) {
	if s.SegmentTimeline == nil || s.SegmentTimeline.SegmentTimelineEntries == nil {
		return nil, errors.New("segment timeline is not defined")
	}

	if s.Timescale == "" {
		return nil, errors.New("timescale is not defined")
	}
	timescale, err := strconv.ParseInt(s.Timescale, 10, 64)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to parse timescale"))
	}

	var result []*SegmentTemplateEntryInfo
	var timestamp int64

	for _, entry := range s.SegmentTimeline.SegmentTimelineEntries {
		repetitions, _ := strconv.ParseInt(entry.RepeatCount, 10, 64)
		repetitions += 1

		duration, err := strconv.ParseInt(entry.Duration, 10, 64)
		if err != nil {
			return nil, err
		}

		for i := int64(0); i < repetitions; i++ {
			result = append(result, &SegmentTemplateEntryInfo{
				TimestampMs: timestamp * 1000 / timescale,
				DurationMs:  duration * 1000 / timescale,
			})
			timestamp += duration
		}
	}

	return result, nil
}

func (s *SegmentTemplate) getChunkDuration() (int64, error) {
	timescale, err := strconv.ParseInt(s.Timescale, 10, 64)
	if err != nil {
		return int64(0), err
	}

	duration, err := s.SegmentTimeline.getChunkDuration()
	if err != nil {
		return int64(0), err
	}

	return duration * 1000 / timescale, nil
}
