package videos

import "regexp"

var (
	initStreamPattern  = regexp.MustCompile(`.*(output-init-stream.*)`)
	chunkStreamPattern = regexp.MustCompile(`.*(output-chunk-stream.*)`)
)

func isChunkStream(filename string) bool {
	return chunkStreamPattern.MatchString(filename)
}

func isInitStream(filename string) bool {
	return initStreamPattern.MatchString(filename)
}
