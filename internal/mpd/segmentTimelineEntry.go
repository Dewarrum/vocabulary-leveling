package mpd

type SegmentTimelineEntry struct {
	Timestamp   string `xml:"t,attr,omitempty" json:"timestamp,omitempty"`
	Duration    string `xml:"d,attr,omitempty" json:"duration,omitempty"`
	RepeatCount string `xml:"r,attr,omitempty" json:"repeatCount,omitempty"`
}
