// Package rsvp is a Go web framework built around content negotiation.
//
// The framework automatically negotiates response format based on the Accept
// header, supporting JSON, XML, HTML, plain text, binary, Gob,
// and MessagePack (using -tags=rsvp_msgpack).
// This content negotiation extends to ALL responses, including redirects,
// allowing you to provide rich feedback in many contexts.
//
// This makes rsvp particularly well-suited for APIs that serve multiple clients
// (browsers, mobile apps, CLI tools) and for taking advantage of principles such
// as REST and progressive enhancement.
package rsvp

import (
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	html "html/template"
	"iter"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	text "text/template"

	"github.com/Teajey/rsvp/internal/content"
	"github.com/Teajey/rsvp/internal/log"
)

type Response struct {
	// Data is the raw data of the response payload to be rendered.
	//
	// IMPORTANT: A nil Data renders as JSON "null\n", not an empty response.
	// Use Data: "" for a blank response body.
	Data any
	// TemplateName sets the template that this Response may attempt to select from
	// [Config.HtmlTemplate] or [Config.TextTemplate],
	//
	// [ResponseWriter.DefaultTemplateName] may also be used to set a default once on a handler.
	//
	// It is not an error if a template is not found for one of the two templates; other formats will be attempted.
	//
	// TODO: Perhaps a warning should be issued to stderr if this fails to match on both templates?
	TemplateName string

	statusCode int

	predeterminedMediaType string

	blankBodyOverride bool

	redirectLocation string
}

// ResponseWriter is a simplified equivalent of [http.ResponseWriter]
type ResponseWriter interface {
	// Header is equivalent to [http.ResponseWriter.Header]
	Header() http.Header

	// DefaultTemplateName is used to associate a default template name with the current handler.
	//
	// It may be overridden by [Response.TemplateName].
	DefaultTemplateName(name string)
}

type response struct {
	std                 http.ResponseWriter
	defaultTemplateName string
}

func (r *response) DefaultTemplateName(name string) {
	r.defaultTemplateName = name
}

func (r *response) Header() http.Header {
	return r.std.Header()
}

func (res *Response) isBlank() bool {
	return res.Data == nil && res.blankBodyOverride
}

var extendedMediaTypes []string = nil

func (res *Response) mediaTypes(cfg Config) iter.Seq[string] {
	return func(yield func(string) bool) {
		if res.predeterminedMediaType != "" {
			log.Dev("Overriding media-types with %s", res.predeterminedMediaType)
			yield(string(res.predeterminedMediaType))
			return
		}

		switch res.Data.(type) {
		case Html:
			yield(SupportedMediaTypeHtml)
			return
		case string:
			if !yield(SupportedMediaTypePlaintext) {
				return
			}
		case []byte:
			yield(SupportedMediaTypeBytes)
			return
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

		if !yield(SupportedMediaTypeGob) {
			return
		}

		for _, mediaType := range extendedMediaTypes {
			if !yield(mediaType) {
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
	}
}

// Settings for writing the rsvp.Response
type Config struct {
	// HtmlTemplate is used by [Response.Write] to potentially render its
	// data to a given HTML template.
	HtmlTemplate *html.Template
	// TextTemplate is used by [Response.Write] to potentially render its
	// data to a given text template.
	TextTemplate *text.Template

	// JsonPrefix is used to set [json.Encoder.SetIndent]
	JsonPrefix string
	// JsonIndent is used to set [json.Encoder.SetIndent]
	JsonIndent string
	// XmlPrefix is used to set [xml.Encoder.Indent]
	XmlPrefix string
	// XmlIndent is used to set [xml.Encoder.Indent]
	XmlIndent string
}

type mediaTypeExtensionHandler = func(mediaType string, w http.ResponseWriter, res *Response) (bool, error)

var mediaTypeExtensionHandlers []mediaTypeExtensionHandler = nil

func (res *Response) determineSupported(cfg Config) []string {
	supported := slices.Collect(res.mediaTypes(cfg))
	log.Dev("supported %v", supported)

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
	log.Dev("mediaType %#v", mediaType)

	return mediaType
}

func (res *Response) determineContentType(mediaType string, h http.Header) {
	contentType := mediaTypeToContentType[mediaType]

	log.Dev("Setting content-type to %#v", contentType)
	h.Set("Content-Type", contentType)
}

// Write the [Response] to the [http.ResponseWriter] with the given [Config].
func (res *Response) Write(w http.ResponseWriter, r *http.Request, cfg Config) error {
	log.Dev("config: %#v", cfg)

	h := w.Header()

	accept := r.Header.Get("Accept")

	contentType := h.Get("Content-Type")
	if contentType != "" {
		aMediaType := string(contentTypeExtractMediaType(contentType))
		_, ok := mediaTypeToContentType[aMediaType]
		if ok {
			res.predeterminedMediaType = aMediaType
			log.Dev("Content-Type is set to a recognised type, so predeterminedMediaType set to %#v", res.predeterminedMediaType)
		}
	}

	ext := determineExt(r)
	supported := res.determineSupported(cfg)
	mediaType := res.determineMediaType(ext, accept, supported)

	if 300 <= res.statusCode && res.statusCode < 400 {
		log.Dev("Redirect")

		h.Set("Location", res.redirectLocation)

		if res.isBlank() || accept == "" {
			log.Dev("Redirect returning empty")
			w.WriteHeader(res.statusCode)
			return nil
		}

		supported := res.determineSupported(cfg)
		mediaType := res.determineMediaType(ext, accept, supported)

		res.determineContentType(mediaType, h)

		w.WriteHeader(res.statusCode)

		return render(res, mediaType, w, cfg)
	}

	if mediaType == "" {
		if ext != "" {
			w.WriteHeader(http.StatusNotFound)
			return nil
		}

		res.statusCode = http.StatusNotAcceptable
		log.Dev("NotAcceptable. Ignoring Accept header...")
		mediaType = chooseMediaType(ext, supported, content.ParseAccept(""))
		log.Dev("new mediaType %#v", mediaType)
	}

	if mediaType == "text/plain" && cfg.TextTemplate == nil && res.TemplateName != "" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	if !res.isBlank() && contentType == "" {
		res.determineContentType(mediaType, h)
	}

	if res.statusCode != 0 {
		w.WriteHeader(res.statusCode)
	}

	if res.isBlank() {
		log.Dev("Early returning because body is empty")
		return nil
	}

	return render(res, mediaType, w, cfg)
}

var ErrFailedToMatchTextTemplate = errors.New("TemplateName was set, but it failed to match within TextTemplate")
var ErrFailedToMatchHtmlTemplate = errors.New("TemplateName was set, but it failed to match within HtmlTemplate")

func render(res *Response, mediaType string, w http.ResponseWriter, cfg Config) error {
	switch mediaType {
	case SupportedMediaTypeHtml:
		log.Dev("Rendering html...")
		if res.TemplateName != "" {
			log.Dev("Template name is set, so expecting a template...")

			if tm := cfg.HtmlTemplate.Lookup(res.TemplateName); tm != nil {
				log.Dev("Executing HtmlTemplate...")
				err := tm.ExecuteTemplate(w, res.TemplateName, res.Data)
				if err != nil {
					return fmt.Errorf("failed to render data in html template %s: %w", res.TemplateName, err)
				}
				break
			}

			return ErrFailedToMatchHtmlTemplate
		}

		_, err := w.Write([]byte(res.Data.(Html)))
		if err != nil {
			return fmt.Errorf("failed to write string as HTML: %w", err)
		}
	case SupportedMediaTypePlaintext:
		log.Dev("Rendering plain text...")

		if res.TemplateName != "" {
			log.Dev("Template name is set, so expecting a template...")

			if tm := cfg.TextTemplate.Lookup(res.TemplateName); tm != nil {
				log.Dev("Executing TextTemplate...")
				err := tm.ExecuteTemplate(w, res.TemplateName, res.Data)
				if err != nil {
					return fmt.Errorf("failed to render data in text template %s: %w", res.TemplateName, err)
				}
				break
			}

			return ErrFailedToMatchTextTemplate
		}
		log.Dev("Not using a template because either TextTemplate or TemplateName is not set...")

		if data, ok := res.Data.(string); ok {
			log.Dev("Can write data directly because it is a string...")
			_, err := w.Write([]byte(data))
			if err != nil {
				return fmt.Errorf("failed to render data as string: %w", err)
			}
			break
		}

		return fmt.Errorf("trying to render data as %s but this type is not supported: %#v", SupportedMediaTypePlaintext, res.Data)
	case SupportedMediaTypeJson:
		log.Dev("Rendering json...")
		enc := json.NewEncoder(w)
		enc.SetIndent(cfg.JsonPrefix, cfg.JsonIndent)
		err := enc.Encode(res.Data)
		if err != nil {
			return fmt.Errorf("failed to render data as JSON: %w", err)
		}
	case SupportedMediaTypeXml:
		log.Dev("Rendering xml...")
		enc := xml.NewEncoder(w)
		enc.Indent(cfg.XmlPrefix, cfg.XmlIndent)
		err := enc.Encode(res.Data)
		if err != nil {
			return fmt.Errorf("failed to render data as XML: %w", err)
		}
	case SupportedMediaTypeCsv:
		data, ok := res.Data.(Csv)
		if !ok {
			return fmt.Errorf("trying to write %#v, but it does not implement rsvp.Csv", res.Data)
		}
		log.Dev("Rendering csv...")
		wr := csv.NewWriter(w)
		err := data.MarshalCsv(wr)
		if err != nil {
			return fmt.Errorf("failed to render data as CSV: %w", err)
		}
	case SupportedMediaTypeBytes:
		log.Dev("Rendering bytes...")
		_, err := w.Write(res.Data.([]byte))
		if err != nil {
			return fmt.Errorf("failed to render data as bytes: %w", err)
		}
	case SupportedMediaTypeGob:
		log.Dev("Rendering gob...")
		err := gob.NewEncoder(w).Encode(res.Data)
		if err != nil {
			return fmt.Errorf("failed to render data as encoding/gob: %w", err)
		}
	default:
		for _, handler := range mediaTypeExtensionHandlers {
			matched, err := handler(mediaType, w, res)
			if err != nil {
				return fmt.Errorf("an extension handler failed: %w", err)
			}
			if matched {
				return nil
			}
		}
		return fmt.Errorf("unhandled mediaType: %#v", mediaType)
	}

	return nil
}

// Blank will render as a blank response with no Content-Type.
//
// Status 200 by default.
func Blank() Response {
	return Response{blankBodyOverride: true}
}

// Data is a convenience function equivalent to instantiating Response{Data: data}
func Data(data any) Response {
	return Response{Data: data}
}

// Html can be used to set [Response.Data].
//
// The wrapped string will be treated as text/html
// instead of text/plain.
type Html string
