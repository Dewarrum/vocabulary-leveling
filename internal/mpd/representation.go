package mpd

type Representation struct {
	ID                        string                     `xml:"id,attr" json:"id,omitempty"`
	Width                     string                     `xml:"width,attr" json:"width,omitempty"`
	Height                    string                     `xml:"height,attr" json:"height,omitempty"`
	MimeType                  string                     `xml:"mimeType,attr" json:"mimeType,omitempty"`
	FrameRate                 string                     `xml:"frameRate,attr" json:"frameRate,omitempty"`
	Bandwidth                 string                     `xml:"bandwidth,attr" json:"bandwidth,omitempty"`
	AudioSamplingRate         string                     `xml:"audioSamplingRate,attr" json:"audioSamplingRate,omitempty"`
	Codecs                    string                     `xml:"codecs,attr" json:"codecs,omitempty"`
	SAR                       string                     `xml:"sar,attr" json:"sar,omitempty"`
	ScanType                  string                     `xml:"scanType,attr" json:"scanType,omitempty"`
	SegmentList               *SegmentList               `xml:"SegmentList,omitempty" json:"segmentList,omitempty"`
	SegmentTemplate           *SegmentTemplate           `xml:"SegmentTemplate,omitempty" json:"segmentTemplate,omitempty"`
	BaseUrl                   string                     `xml:"BaseURL,omitempty" json:"baseUrl,omitempty"`
	AudioChannelConfiguration *AudioChannelConfiguration `xml:"AudioChannelConfiguration,omitempty" json:"audioChannelConfiguration,omitempty"`
}

type AudioChannelConfiguration struct {
	SchemeIdUri string `xml:"schemeIdUri,attr"`
	Value       string `xml:"value,attr"`
}
