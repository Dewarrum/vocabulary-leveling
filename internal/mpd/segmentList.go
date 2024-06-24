package mpd

type SegmentList struct {
	Timescale      string          `xml:"timescale,attr"`
	Duration       string          `xml:"duration,attr"`
	StartNumber    string          `xml:"startNumber,attr"`
	Initialization *Initialization `xml:"Initialization"`
	Segments       []*Segment      `xml:"SegmentURL"`
}
