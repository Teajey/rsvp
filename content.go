package rsvp

import (
	"iter"
	"strings"

	"github.com/Teajey/rsvp/internal/log"
)

const (
	SupportedMediaTypePlaintext string = "text/plain"
	SupportedMediaTypeHtml      string = "text/html"
	SupportedMediaTypeBytes     string = "application/octet-stream"
	SupportedMediaTypeJson      string = "application/json"
	SupportedMediaTypeXml       string = "application/xml"
	SupportedMediaTypeGob       string = "application/vnd.golang.gob"
)

// Can be used to place every supported type on:
//   - Config.RenderSeeOtherBlackList
//   - Config.RenderFoundBlackList
var AllSupportedTypes = []string{
	SupportedMediaTypePlaintext,
	SupportedMediaTypeHtml,
	SupportedMediaTypeBytes,
	SupportedMediaTypeJson,
	SupportedMediaTypeXml,
	SupportedMediaTypeGob,
}

var mediaTypeToContentType = map[string]string{
	// TODO: Why did I insist on specifying utf-8 here? There should be note. I think it might just be because it's inline with what net/http does
	SupportedMediaTypePlaintext: "text/plain; charset=utf-8",
	SupportedMediaTypeHtml:      "text/html; charset=utf-8",
	SupportedMediaTypeBytes:     "application/octet-stream",
	SupportedMediaTypeJson:      "application/json",
	SupportedMediaTypeXml:       "application/xml",
	SupportedMediaTypeGob:       "application/vnd.golang.gob",
}

var extToProposalMap = map[string]string{
	"txt":  SupportedMediaTypePlaintext,
	"html": SupportedMediaTypeHtml,
	"htm":  SupportedMediaTypeHtml,
	"json": SupportedMediaTypeJson,
	"xml":  SupportedMediaTypeXml,
	"bin":  SupportedMediaTypeBytes,
	"gob":  SupportedMediaTypeGob,
}

// mediatype string m must be well-formed
func splitMediaType(m string) (string, string) {
	sp := strings.SplitN(m, "/", 2)
	return sp[0], sp[1]
}

// mediatype string b must be well-formed
func mediaTypesEqual(a string, b string) bool {
	if strings.HasPrefix(b, "*/") {
		return true
	}

	b1, b2 := splitMediaType(b)
	s := a
	s1, s2 := splitMediaType(s)

	if b1 == s1 {
		return b2 == "*" || b2 == s2
	}

	return false
}

func chooseMediaType(ext string, supported []string, accept iter.Seq[string]) string {
	if ext != "" {
		log.Dev("Checking extension: %#v", ext)
		if a, ok := extToProposalMap[ext]; ok {
			log.Dev("proposing %#v", a)
			for _, s := range supported {
				log.Dev("offering  %#v", s)
				if mediaTypesEqual(s, a) {
					return s
				}
			}
		}
		return ""
	}

	log.Dev("Checking accept list")
	for a := range accept {
		log.Dev("proposing %#v", a)
		for _, s := range supported {
			log.Dev("offering  %#v", s)
			if mediaTypesEqual(s, a) {
				return s
			}
		}
	}

	return ""
}
