package mpd

type AdaptationSet struct {
	Id                      string            `xml:"id,attr"`
	MimeType                string            `xml:"mimeType,attr"`
	ContentType             string            `xml:"contentType,attr"`
	SegmentAlignment        string            `xml:"segmentAlignment,attr"`
	SubsegmentAlignment     string            `xml:"subsegmentAlignment,attr"`
	StartWithSAP            string            `xml:"startWithSAP,attr"`
	MaxWidth                *string           `xml:"maxWidth,attr"`
	MaxHeight               *string           `xml:"maxHeight,attr"`
	SubsegmentStartsWithSAP string            `xml:"subsegmentStartsWithSAP,attr"`
	BitstreamSwitching      *bool             `xml:"bitstreamSwitching,attr"`
	Lang                    *string           `xml:"lang,attr"`
	Par                     *string           `xml:"par,attr"`
	Codecs                  *string           `xml:"codecs,attr"`
	Representations         []*Representation `xml:"Representation,omitempty"`
}
