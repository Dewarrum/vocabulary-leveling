package mpd

import "strconv"

type SegmentTimeline struct {
	SegmentTimelineEntries []*SegmentTimelineEntry `xml:"S,omitempty"`
}

func (s *SegmentTimeline) GetTotalDuration() (int64, error) {
	var totalDuration int64
	for _, entry := range s.SegmentTimelineEntries {
		repetitions, _ := strconv.ParseInt(entry.RepeatCount, 10, 64)
		repetitions += 1

		duration, err := strconv.ParseInt(entry.Duration, 10, 64)
		if err != nil {
			return int64(0), err
		}

		for i := int64(0); i < repetitions; i++ {
			totalDuration += duration
		}
	}

	return totalDuration, nil
}

func (s *SegmentTimeline) getChunkDuration() (int64, error) {
	var chunkDuration int64
	for _, entry := range s.SegmentTimelineEntries {
		duration, err := strconv.ParseInt(entry.Duration, 10, 64)
		if err != nil {
			return int64(0), err
		}
		if duration > chunkDuration {
			chunkDuration = duration
		}
	}

	return chunkDuration, nil
}
