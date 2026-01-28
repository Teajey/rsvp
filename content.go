package rsvp

import (
	"io"
	"iter"
	"net/http"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Teajey/rsvp/internal/content"
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

var extendedMediaTypes []string = nil

type mediaTypeExtensionHandler = func(mediaType string, w io.Writer, res *Response) (bool, error)

var mediaTypeExtensionHandlers []mediaTypeExtensionHandler = nil

func (res *Response) determineSupported(cfg Config) []string {
	supported := slices.Collect(res.MediaTypes(cfg))
	dev.Log("supported %v", supported)

	return supported
}

func determineExt(r *http.Request) string {
	var ext string
	if r.Method == http.MethodGet {
		ext = strings.TrimPrefix(filepath.Ext(r.URL.Path), ".")
	}

	return ext
}

func (res *Response) determineMediaType(ext, accept string, supported []string) string {
	mediaType := chooseMediaType(ext, supported, content.ParseAccept(accept))
	dev.Log("mediaType %#v", mediaType)

	return mediaType
}

func (res *Response) determineContentType(mediaType string, wh http.Header) {
	contentType := mediaTypeToContentType[mediaType]

	dev.Log("Setting content-type to %#v", contentType)
	wh.Set("Content-Type", contentType)
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

// MediaTypes returns the sequence of media types (e.g. text/plain) in the order that this [Response] will propose.
//
// The order generally follows this pattern:
//  1. Type-specific (Html wrapper, string, bytes)
//  2. Generic structured (JSON, XML)
//  3. Interface implementations (CSV)
//  4. Template-based (HTML template, text template)
//  5. Golang-native fallback (Gob)
func (res *Response) MediaTypes(cfg Config) iter.Seq[string] {
	return func(yield func(string) bool) {
		if res.predeterminedMediaType != "" {
			dev.Log("Overriding media-types with %s", res.predeterminedMediaType)
			yield(string(res.predeterminedMediaType))
			return
		}

		switch res.Data.(type) {
		case string:
			if !yield(SupportedMediaTypePlaintext) {
				return
			}
		case []byte:
			if !yield(SupportedMediaTypeBytes) {
				return
			}
		}

		if !yield(SupportedMediaTypeJson) {
			return
		}

		if !yield(SupportedMediaTypeXml) {
			return
		}

		_, ok := res.Data.(Csv)
		if ok {
			if !yield(SupportedMediaTypeCsv) {
				return
			}
		}

		if res.TemplateName != "" {
			if cfg.HtmlTemplate != nil {
				if !yield(SupportedMediaTypeHtml) {
					return
				}
			}

			if cfg.TextTemplate != nil {
				if !yield(SupportedMediaTypePlaintext) {
					return
				}
			}
		}

		if !yield(SupportedMediaTypeGob) {
			return
		}

		for _, mediaType := range extendedMediaTypes {
			if !yield(mediaType) {
				return
			}
		}
	}
}
