package videos_test

import (
	"dewarrum/vocabulary-leveling/internal/videos"
	"testing"
)

func TestExtendRangeWhenRangeIsBigger(t *testing.T) {
	startMs := int64(1000)
	endMs := int64(5000)
	desiredDuration := int64(3000)

	s, e := videos.ExtendRange(startMs, endMs, desiredDuration)

	if s != startMs {
		t.Errorf("Expected startMs to be %d, but got %d", startMs, s)
	}

	if e != endMs {
		t.Errorf("Expected endMs to be %d, but got %d", endMs, e)
	}
}

func TestExtendRangeWhenRangeIsSmaller(t *testing.T) {
	startMs := int64(1000)
	endMs := int64(2000)
	desiredDuration := int64(2000)

	s, e := videos.ExtendRange(startMs, endMs, desiredDuration)

	if s != 500 {
		t.Errorf("Expected startMs to be %d, but got %d", 0, s)
	}

	if e != 2500 {
		t.Errorf("Expected endMs to be %d, but got %d", endMs, e)
	}
}

func TestExtendRangeWhenStartIsZero(t *testing.T) {
	startMs := int64(0)
	endMs := int64(1000)
	desiredDuration := int64(2000)

	s, e := videos.ExtendRange(startMs, endMs, desiredDuration)

	if s != 0 {
		t.Errorf("Expected startMs to be %d, but got %d", 0, s)
	}

	if e != 2000 {
		t.Errorf("Expected endMs to be %d, but got %d", 2000, e)
	}
}
