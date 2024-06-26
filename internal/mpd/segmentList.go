package mpd

type SegmentList struct {
	Timescale      string          `xml:"timescale,attr,omitempty"`
	Duration       string          `xml:"duration,attr,omitempty"`
	StartNumber    string          `xml:"startNumber,attr,omitempty"`
	Initialization *Initialization `xml:"Initialization,omitempty"`
	Segments       []*Segment      `xml:"SegmentURL,omitempty"`
}
