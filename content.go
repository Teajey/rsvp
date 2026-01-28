package rsvp

import (
	"iter"
	"strings"

	"github.com/Teajey/rsvp/internal/dev"
)

const (
	SupportedMediaTypePlaintext string = "text/plain"
	SupportedMediaTypeHtml      string = "text/html"
	SupportedMediaTypeCsv       string = "text/csv"
	SupportedMediaTypeBytes     string = "application/octet-stream"
	SupportedMediaTypeJson      string = "application/json"
	SupportedMediaTypeXml       string = "application/xml"
	SupportedMediaTypeGob       string = "application/vnd.golang.gob"
)

var mediaTypeToContentType = map[string]string{
	// TODO: Why did I insist on specifying utf-8 here? There should be note. I think it might just be because it's inline with what net/http does
	SupportedMediaTypePlaintext: "text/plain; charset=utf-8",
	SupportedMediaTypeHtml:      "text/html; charset=utf-8",
	SupportedMediaTypeCsv:       "text/csv; charset=utf-8",
	SupportedMediaTypeBytes:     "application/octet-stream",
	SupportedMediaTypeJson:      "application/json",
	SupportedMediaTypeXml:       "application/xml",
	SupportedMediaTypeGob:       "application/vnd.golang.gob",
}

var extToProposalMap = map[string]string{
	"txt":  SupportedMediaTypePlaintext,
	"html": SupportedMediaTypeHtml,
	"htm":  SupportedMediaTypeHtml,
	"csv":  SupportedMediaTypeCsv,
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
		dev.Log("Checking extension: %#v", ext)
		if a, ok := extToProposalMap[ext]; ok {
			dev.Log("proposing %#v", a)
			for _, s := range supported {
				dev.Log("offering  %#v", s)
				if mediaTypesEqual(s, a) {
					return s
				}
			}
		}
		return ""
	}

	dev.Log("Checking accept list")
	for a := range accept {
		dev.Log("proposing %#v", a)
		for _, s := range supported {
			dev.Log("offering  %#v", s)
			if mediaTypesEqual(s, a) {
				return s
			}
		}
	}

	return ""
}
