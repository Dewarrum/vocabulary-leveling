package manifests

type DbPeriod struct {
	Id             string             `json:"id"`
	Start          string             `json:"start"`
	AdaptationSets []*DbAdaptationSet `json:"adaptationSets"`
}

type DbAdaptationSet struct {
	Id                 string              `json:"id"`
	ContentType        string              `json:"contentType"`
	StartWithSap       string              `json:"startWithSap"`
	SegmentAlignment   string              `json:"segmentAlignment"`
	BitstreamSwitching *bool               `json:"bitstreamSwitching"`
	Lang               *string             `json:"lang"`
	MaxWidth           *string             `json:"maxWidth"`
	MaxHeight          *string             `json:"maxHeight"`
	Par                *string             `json:"par"`
	Representations    []*DbRepresentation `json:"representations"`
}
