package mpd

type SegmentList struct {
	Timescale      string          `xml:"timescale,attr,omitempty" json:"timescale,omitempty"`
	Duration       string          `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	StartNumber    string          `xml:"startNumber,attr,omitempty" json:"startNumber,omitempty"`
	Initialization *Initialization `xml:"Initialization,omitempty" json:"initialization,omitempty"`
	Segments       []*Segment      `xml:"SegmentURL,omitempty" json:"segments,omitempty"`
}
