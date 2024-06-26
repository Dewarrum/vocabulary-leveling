package mpd

type Period struct {
	Start          string           `xml:"start,attr"`
	ID             string           `xml:"id,attr"`
	Duration       string           `xml:"duration,attr,omitempty"`
	AdaptationSets []*AdaptationSet `xml:"AdaptationSet"`
}
