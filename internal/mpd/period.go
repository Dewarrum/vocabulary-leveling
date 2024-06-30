package mpd

type Period struct {
	Start          string           `xml:"start,attr" json:"start,omitempty"`
	ID             string           `xml:"id,attr" json:"id,omitempty"`
	Duration       string           `xml:"duration,attr,omitempty" json:"duration,omitempty"`
	AdaptationSets []*AdaptationSet `xml:"AdaptationSet" json:"adaptationSets,omitempty"`
}

func (p *Period) getChunkDuration() (int64, error) {
	var chunkDuration int64
	for _, adaptationSet := range p.AdaptationSets {
		for _, representation := range adaptationSet.Representations {
			duration, err := representation.SegmentTemplate.getChunkDuration()
			if err != nil {
				return int64(0), err
			}
			if duration > chunkDuration {
				chunkDuration = duration
			}
		}
	}

	return chunkDuration, nil
}
