package rsvp

import (
	"iter"
	"net/url"
	"path/filepath"
	"strings"
)

type supportedType string

const (
	mPlaintext supportedType = "text/plain"
	mHtml      supportedType = "text/html"
	mBytes     supportedType = "application/octet-stream"
	mJson      supportedType = "application/json"
)

var ext2mediatype = map[string]string{
	".md":   "*/*",
	".txt":  string(mPlaintext),
	".html": string(mHtml),
	".json": string(mJson),
}

// mediatype string m must be well-formed
func splitMediaType(m string) (string, string) {
	sp := strings.SplitN(m, "/", 2)
	return sp[0], sp[1]
}

// mediatype string b must be well-formed
func mediaTypesEqual(a supportedType, b string) bool {
	if strings.HasPrefix(b, "*/") {
		return true
	}

	b1, b2 := splitMediaType(b)
	s := string(a)
	s1, s2 := splitMediaType(s)

	if b1 == s1 {
		return b2 == "*" || b2 == s2
	}

	return false
}

func resolveMediaType(u *url.URL, supported []supportedType, accept iter.Seq[string]) string {
	ext := filepath.Ext(u.Path)

	if ext != "" {
		if m, ok := ext2mediatype[ext]; ok {
			for _, s := range supported {
				if mediaTypesEqual(s, m) {
					return string(s)
				}
			}
			return ""
		}
	}

	for a := range accept {
		for _, s := range supported {
			if mediaTypesEqual(s, a) {
				return string(s)
			}
		}
	}

	return ""
}
