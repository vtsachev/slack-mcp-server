package text

import (
	"strings"

	"github.com/bbalet/stopwords"
)

func ProcessText(s string) string {
	s = stopwordsFilter(s)
	s = strings.TrimSpace(s)
	return s
}

func stopwordsFilter(s string) string {
	return stopwords.CleanString(s, "en", true)
}
