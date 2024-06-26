package mpd

type SegmentTemplate struct {
	Timescale       string          `xml:"timescale,attr"`
	Initialization  *Initialization `xml:"Initialization"`
	Media           string          `xml:"media,attr"`
	Duration        string          `xml:"duration,attr"`
	StartNumber     string          `xml:"startNumber,attr"`
	Times           string          `xml:"times,attr"`
	Presentation    string          `xml:"presentation,attr"`
	Bandwidth       string          `xml:"bandwidth,attr"`
	ProgramDateTime string          `xml:"programDateTime,attr"`
	SegmentTimeline *SegmentTimeline
}
