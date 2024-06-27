package mpd_test

import (
	"dewarrum/vocabulary-leveling/internal/mpd"
	"testing"
)

func TestGetSegmentInfosSingleEntryWithoutRepetitions(t *testing.T) {
	segmentTemplate := &mpd.SegmentTemplate{
		Timescale: "1000",
		SegmentTimeline: &mpd.SegmentTimeline{
			SegmentTimelineEntries: []*mpd.SegmentTimelineEntry{
				{
					Duration: "20000",
				},
			},
		},
	}

	segmentInfos, err := segmentTemplate.GetSegmentInfos()
	if err != nil {
		t.Error(err)
	}

	if len(segmentInfos) != 1 {
		t.Error("Expected 1 segment info")
	}

	expectedDuration := int64(20000)
	if segmentInfos[0].DurationMs != expectedDuration {
		t.Errorf("Expected duration to be %d, but got %d", expectedDuration, segmentInfos[0].DurationMs)
	}

	expectedTimestamp := int64(0)
	if segmentInfos[0].TimestampMs != 0 {
		t.Errorf("Expected timestamp to be %d, but got %d", expectedTimestamp, segmentInfos[0].TimestampMs)
	}
}

func TestGetSegmentInfosMultipleEntries(t *testing.T) {
	segmentTemplate := &mpd.SegmentTemplate{
		Timescale: "1000",
		SegmentTimeline: &mpd.SegmentTimeline{
			SegmentTimelineEntries: []*mpd.SegmentTimelineEntry{
				{
					Duration: "20000",
				},
				{
					Duration: "20000",
				},
			},
		},
	}

	segmentInfos, err := segmentTemplate.GetSegmentInfos()
	if err != nil {
		t.Error(err)
	}

	if len(segmentInfos) != 2 {
		t.Error("Expected 2 segment info")
	}

	expectedDuration := int64(20000)
	if segmentInfos[0].DurationMs != expectedDuration {
		t.Errorf("Expected duration to be %d, but got %d", expectedDuration, segmentInfos[0].DurationMs)
	}

	expectedTimestamp := int64(0)
	if segmentInfos[0].TimestampMs != expectedTimestamp {
		t.Errorf("Expected timestamp to be %d, but got %d", expectedTimestamp, segmentInfos[0].TimestampMs)
	}

	expectedDuration = int64(20000)
	if segmentInfos[1].DurationMs != expectedDuration {
		t.Errorf("Expected duration to be %d, but got %d", expectedDuration, segmentInfos[1].DurationMs)
	}

	expectedTimestamp = int64(20000)
	if segmentInfos[1].TimestampMs != expectedTimestamp {
		t.Errorf("Expected timestamp to be %d, but got %d", expectedTimestamp, segmentInfos[1].TimestampMs)
	}
}

func TestGetSegmentInfosSingleEntryWithRepetitions(t *testing.T) {
	segmentTemplate := &mpd.SegmentTemplate{
		Timescale: "1000",
		SegmentTimeline: &mpd.SegmentTimeline{
			SegmentTimelineEntries: []*mpd.SegmentTimelineEntry{
				{
					Duration:    "20000",
					RepeatCount: "1",
				},
			},
		},
	}

	segmentInfos, err := segmentTemplate.GetSegmentInfos()
	if err != nil {
		t.Error(err)
	}

	if len(segmentInfos) != 2 {
		t.Error("Expected 2 segment info")
	}

	expectedDuration := int64(20000)
	if segmentInfos[0].DurationMs != expectedDuration {
		t.Errorf("Expected duration to be %d, but got %d", expectedDuration, segmentInfos[0].DurationMs)
	}

	expectedTimestamp := int64(0)
	if segmentInfos[0].TimestampMs != expectedTimestamp {
		t.Errorf("Expected timestamp to be %d, but got %d", expectedTimestamp, segmentInfos[0].TimestampMs)
	}

	expectedDuration = int64(20000)
	if segmentInfos[1].DurationMs != expectedDuration {
		t.Errorf("Expected duration to be %d, but got %d", expectedDuration, segmentInfos[1].DurationMs)
	}

	expectedTimestamp = int64(20000)
	if segmentInfos[1].TimestampMs != expectedTimestamp {
		t.Errorf("Expected timestamp to be %d, but got %d", expectedTimestamp, segmentInfos[1].TimestampMs)
	}
}

func TestGetSegmentInfosMultipleAndSingleRepetitions(t *testing.T) {
	segmentTemplate := &mpd.SegmentTemplate{
		Timescale: "1000",
		SegmentTimeline: &mpd.SegmentTimeline{
			SegmentTimelineEntries: []*mpd.SegmentTimelineEntry{
				{
					Duration:    "20000",
					RepeatCount: "2",
				},
				{
					Duration:    "20000",
					RepeatCount: "1",
				},
			},
		},
	}

	segmentInfos, err := segmentTemplate.GetSegmentInfos()
	if err != nil {
		t.Error(err)
	}

	if len(segmentInfos) != 5 {
		t.Error("Expected 5 segment info, but got", len(segmentInfos))
	}

	lastSegmentInfo := segmentInfos[4]
	expectedDuration := int64(20000)
	if lastSegmentInfo.DurationMs != expectedDuration {
		t.Errorf("Expected duration to be %d, but got %d", expectedDuration, lastSegmentInfo.DurationMs)
	}

	expectedTimestamp := int64(80000)
	if lastSegmentInfo.TimestampMs != expectedTimestamp {
		t.Errorf("Expected timestamp to be %d, but got %d", expectedTimestamp, lastSegmentInfo.TimestampMs)
	}
}

func TestGetSegmentInfosMultipleAndMissingRepetitions(t *testing.T) {
	segmentTemplate := &mpd.SegmentTemplate{
		Timescale: "1000",
		SegmentTimeline: &mpd.SegmentTimeline{
			SegmentTimelineEntries: []*mpd.SegmentTimelineEntry{
				{
					Duration:    "20000",
					RepeatCount: "1",
				},
				{
					Duration: "20000",
				},
			},
		},
	}

	segmentInfos, err := segmentTemplate.GetSegmentInfos()
	if err != nil {
		t.Error(err)
	}

	if len(segmentInfos) != 3 {
		t.Error("Expected 3 segment info, but got", len(segmentInfos))
	}

	segmentInfo := segmentInfos[2]
	expectedDuration := int64(20000)
	if segmentInfo.DurationMs != expectedDuration {
		t.Errorf("Expected duration to be %d, but got %d", expectedDuration, segmentInfo.DurationMs)
	}

	expectedTimestamp := int64(40000)
	if segmentInfo.TimestampMs != expectedTimestamp {
		t.Errorf("Expected timestamp to be %d, but got %d", expectedTimestamp, segmentInfo.TimestampMs)
	}
}
