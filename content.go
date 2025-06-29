package rsvp

import (
	"iter"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Teajey/rsvp/internal/log"
)

// TODO: I'm not sure that this type is helpful. It might be worth removing.
type supportedType string

const (
	mPlaintext supportedType = "text/plain"
	mHtml      supportedType = "text/html"
	mBytes     supportedType = "application/octet-stream"
	mJson      supportedType = "application/json"
)

var mediaTypeToContentType = map[supportedType]string{
	mPlaintext: "text/plain; charset=utf-8",
	mHtml:      "text/html; charset=utf-8",
	mBytes:     "application/octet-stream",
	mJson:      "application/json",
}

var defaultExtToProposalMap = map[string]string{
	"txt":  string(mPlaintext),
	"html": string(mHtml),
	"json": string(mJson),
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

func chooseMediaType(u *url.URL, supported []supportedType, accept iter.Seq[string], ext2proposal map[string]string) supportedType {
	ext := strings.TrimPrefix(filepath.Ext(u.Path), ".")

	if ext != "" {
		log.Dev("Checking extension: %#v", ext)
		if a, ok := ext2proposal[ext]; ok {
			log.Dev("a %#v", a)
			for _, s := range supported {
				log.Dev("s %#v", s)
				if mediaTypesEqual(s, a) {
					return s
				}
			}
			return ""
		}
	}

	log.Dev("Checking accept header")
	for a := range accept {
		log.Dev("a %#v", a)
		for _, s := range supported {
			log.Dev("s %#v", s)
			if mediaTypesEqual(s, a) {
				return s
			}
		}
	}

	return ""
}
