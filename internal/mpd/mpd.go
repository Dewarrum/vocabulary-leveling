package mpd

import "encoding/xml"

type MPD struct {
	XMLNS                     string    `xml:"xmlns,attr"`
	Type                      string    `xml:"type,attr"`
	MediaPresentationDuration string    `xml:"mediaPresentationDuration,attr"`
	MinBufferTime             string    `xml:"minBufferTime,attr"`
	Profiles                  string    `xml:"profiles,attr"`
	Periods                   []*Period `xml:"Period"`
}

func Parse(b []byte) (*MPD, error) {
	m := new(MPD)
	err := xml.Unmarshal(b, m)
	return m, err
}

func (m *MPD) Serialize() ([]byte, error) {
	return xml.MarshalIndent(m, "", "  ")
}
