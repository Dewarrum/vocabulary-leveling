package mpd

type Period struct {
	Start          string           `xml:"start,attr" json:"start,omitempty"`
	ID             string           `xml:"id,attr" json:"id,omitempty"`
	Duration       string           `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	AdaptationSets []*AdaptationSet `xml:"AdaptationSet" json:"adaptationSets,omitempty"`
}
