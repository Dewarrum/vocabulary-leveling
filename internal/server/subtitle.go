package server

type DtoSubtitle struct {
	Id        string `json:"id"`
	VideoName string `json:"videoName"`
	StartMs   int64  `json:"startMs"`
	EndMs     int64  `json:"endMs"`
	Text      string `json:"text"`
}
