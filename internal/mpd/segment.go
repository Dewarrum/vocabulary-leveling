package mpd

type Segment struct {
	Media string `xml:"media,attr" json:"media,omitempty"`
}
