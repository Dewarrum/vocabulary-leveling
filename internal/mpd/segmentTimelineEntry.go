package mpd

type SegmentTimelineEntry struct {
	Timestamp   string `xml:"t,attr,omitempty"`
	Duration    string `xml:"d,attr,omitempty"`
	RepeatCount string `xml:"r,attr,omitempty"`
}
