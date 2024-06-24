package mpd

type Representation struct {
	ID                string       `xml:"id,attr"`
	Width             string       `xml:"width,attr"`
	Height            string       `xml:"height,attr"`
	FrameRate         string       `xml:"frameRate,attr"`
	Bandwidth         string       `xml:"bandwidth,attr"`
	AudioSamplingRate string       `xml:"audioSamplingRate,attr"`
	Codecs            string       `xml:"codecs,attr"`
	SAR               string       `xml:"sar,attr"`
	ScanType          string       `xml:"scanType,attr"`
	SegmentList       *SegmentList `xml:"SegmentList"`
	BaseUrl           string       `xml:"BaseURL,omitempty"`
}
