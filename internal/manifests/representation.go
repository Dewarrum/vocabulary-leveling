package manifests

type DbRepresentation struct {
	Id                        string                       `json:"id"`
	MimeType                  string                       `json:"mimeType"`
	Codecs                    string                       `json:"codecs"`
	Bandwidth                 string                       `json:"bandwidth"`
	AudioSamplingRate         string                       `json:"audioSamplingRate"`
	Width                     *string                      `json:"width"`
	Height                    *string                      `json:"height"`
	Sar                       *string                      `json:"sar"`
	AudioChannelConfiguration *DbAudioChannelConfiguration `json:"audioChannelConfiguration,omitempty"`
	SegmentTemplate           *DbSegmentTemplate           `json:"segmentTemplate"`
}

type DbAudioChannelConfiguration struct {
	SchemeUriId string `json:"schemeUriId"`
	Value       string `json:"value"`
}
