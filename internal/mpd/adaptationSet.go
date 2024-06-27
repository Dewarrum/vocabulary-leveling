package mpd

type AdaptationSet struct {
	Id                      string            `xml:"id,attr" json:"id,omitempty"`
	MimeType                string            `xml:"mimeType,attr" json:"mimeType,omitempty"`
	ContentType             string            `xml:"contentType,attr" json:"contentType,omitempty"`
	SegmentAlignment        string            `xml:"segmentAlignment,attr" json:"segmentAlignment,omitempty"`
	SubsegmentAlignment     string            `xml:"subsegmentAlignment,attr" json:"subsegmentAlignment,omitempty"`
	StartWithSAP            string            `xml:"startWithSAP,attr" json:"startWithSap,omitempty"`
	MaxWidth                *string           `xml:"maxWidth,attr" json:"maxWidth,omitempty"`
	MaxHeight               *string           `xml:"maxHeight,attr" json:"maxHeight,omitempty"`
	SubsegmentStartsWithSAP string            `xml:"subsegmentStartsWithSAP,attr" json:"subsegmentStartsWithSap,omitempty"`
	BitstreamSwitching      *bool             `xml:"bitstreamSwitching,attr" json:"bitstreamSwitching,omitempty"`
	Lang                    *string           `xml:"lang,attr" json:"lang,omitempty"`
	Par                     *string           `xml:"par,attr" json:"par,omitempty"`
	Codecs                  *string           `xml:"codecs,attr" json:"codecs,omitempty"`
	Representations         []*Representation `xml:"Representation,omitempty" json:"representations,omitempty"`
}
