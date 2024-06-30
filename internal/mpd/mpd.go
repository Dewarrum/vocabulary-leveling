package mpd

import "encoding/xml"

type MPD struct {
	XMLNS                     string    `xml:"xmlns,attr" json:"xmlns,omitempty"`
	Type                      string    `xml:"type,attr" json:"type,omitempty"`
	MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr" json:"mediaPresentationDuration,omitempty"`
	MinBufferTime             string    `xml:"minBufferTime,attr" json:"minBufferTime,omitempty"`
	Profiles                  string    `xml:"profiles,attr" json:"profiles,omitempty"`
	Periods                   []*Period `xml:"Period" json:"periods,omitempty"`
}

func (m *MPD) GetRepresentation(representationId string) *Representation {
	if len(m.Periods) < 1 {
		return nil
	}
	period := m.Periods[0]

	if len(period.AdaptationSets) < 1 {
		return nil
	}

	var adaptationSet *AdaptationSet
	for _, as := range period.AdaptationSets {
		if as.Id == representationId {
			adaptationSet = as
		}
	}

	if adaptationSet == nil {
		return nil
	}

	for _, representation := range adaptationSet.Representations {
		if representation.ID == representationId {
			return representation
		}
	}

	return nil
}

func Parse(b []byte) (*MPD, error) {
	m := new(MPD)
	err := xml.Unmarshal(b, m)
	return m, err
}

func (m *MPD) Serialize() ([]byte, error) {
	return xml.MarshalIndent(m, "", "  ")
}

func (m *MPD) GetChunkDuration() (int64, error) {
	var chunkDuration int64
	for _, period := range m.Periods {
		duration, err := period.getChunkDuration()
		if err != nil {
			return int64(0), err
		}
		if duration > chunkDuration {
			chunkDuration = duration
		}
	}

	return chunkDuration, nil
}
