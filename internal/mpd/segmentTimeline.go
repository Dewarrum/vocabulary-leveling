package mpd

type SegmentTimeline struct {
	SegmentTemplateEntry []*SegmentTimelineEntry `xml:"SegmentTemplateEntry,omitempty"`
}
