package mpd

type Initialization struct {
	SourceURL string `xml:"sourceURL,attr" json:"sourceURL,omitempty"`
}
