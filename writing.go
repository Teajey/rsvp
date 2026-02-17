package rsvp

import (
	"cmp"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/Teajey/rsvp/internal/content"
	"github.com/Teajey/rsvp/internal/dev"
)

// ResponseWriter handles metadata and configuration of the response. It bears its "Writer" name mostly for the sake of keeping rsvp.Handler similar to http.Handler.
//
// Its underlying type has a `write` function, but it is not available here because it is controlled indirectly by the [Body] value that [Handler] provides.
//
// If you need access to http.ResponseWriter, especially for middleware, you should follow the example of [HandleFunc]'s source code for how to operate rsvp at a lower level from within an [http.Handler].
type ResponseWriter interface {
	// Header is equivalent to [http.ResponseWriter.Header]
	Header() http.Header

	// DefaultTemplateName is used to associate a default template name with the current handler.
	//
	// It may be overridden by [Body.TemplateName].
	//
	// The intended use case for this method is to call it at the top of an [HandlerFunc] so that the TemplateName does not need to be set exhaustively on every instance of [Body] that the handler might return.
	DefaultTemplateName(name string)
}

// Write the result of handler to w. May write headers to w.Header().
//
// NOTE: This function is for advanced lower-level use cases.
//
// This function should be used to wrap [Handler] in middleware that requires write access to [http.ResponseWriter].
//
// See this test for an example: https://github.com/Teajey/rsvp/blob/main/middleware_test.go
func Write(w http.ResponseWriter, cfg Config, r *http.Request, handler HandlerFunc) error {
	rw := responseWriter{
		writer: w,
	}
	response := handler(&rw, r)
	return rw.write(&response, r, cfg)
}

type responseWriter struct {
	writer              http.ResponseWriter
	defaultTemplateName string
}

func (w *responseWriter) DefaultTemplateName(name string) {
	w.defaultTemplateName = name
}

func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write the [Body] to the [http.ResponseWriter] with the given [Config].
func (w *responseWriter) write(res *Body, r *http.Request, cfg Config) (err error) {
	dev.Log("config: %#v", cfg)
	status := cmp.Or(res.statusCode, 200)

	if res.TemplateName == "" && w.defaultTemplateName != "" {
		dev.Log("Using default template name: %v", w.defaultTemplateName)
		res.TemplateName = w.defaultTemplateName
	}

	accept := r.Header.Get("Accept")

	wh := w.Header()

	contentType := wh.Get("Content-Type")
	if contentType != "" {
		aMediaType := string(contentTypeExtractMediaType(contentType))
		_, ok := mediaTypeToContentType[aMediaType]
		if ok {
			res.predeterminedMediaType = aMediaType
			dev.Log("Content-Type is set to a recognised type, so predeterminedMediaType set to %#v", res.predeterminedMediaType)
		}
	}

	ext := determineExt(r)
	supported := res.determineSupported(cfg)
	mediaType := res.determineMediaType(ext, accept, supported)

	if 300 <= res.statusCode && res.statusCode < 400 {
		dev.Log("Redirect")

		wh.Set("Location", res.redirectLocation)

		if res.isBlank() || accept == "" {
			dev.Log("Redirect returning empty")
			w.writer.WriteHeader(status)
			return
		}

		supported := res.determineSupported(cfg)
		mediaType := res.determineMediaType(ext, accept, supported)

		res.determineContentType(mediaType, wh)

		w.writer.WriteHeader(status)
		err = render(res, mediaType, w.writer, cfg)
		return
	}

	if mediaType == "" {
		if ext != "" {
			a, ok := extToProposalMap[ext]
			if !ok || !slices.Contains(supported, a) {
				w.writer.WriteHeader(http.StatusNotFound)
				return
			}
		}

		dev.Log("NotAcceptable. Ignoring Accept header and setting status code to 406...")
		status = http.StatusNotAcceptable
		mediaType = chooseMediaType(ext, supported, content.ParseAccept(""))
		dev.Log("new mediaType %#v", mediaType)
	}

	if mediaType == "text/plain" && cfg.TextTemplate == nil && res.TemplateName != "" {
		w.writer.WriteHeader(http.StatusNotFound)
		return
	}

	if !res.isBlank() && contentType == "" {
		res.determineContentType(mediaType, wh)
	}

	if res.isBlank() {
		dev.Log("Early returning because body is empty")
		w.writer.WriteHeader(status)
		return
	}

	w.writer.WriteHeader(status)
	err = render(res, mediaType, w.writer, cfg)
	return
}

var ErrFailedToMatchTextTemplate = errors.New("TemplateName was set, but it failed to match within TextTemplate")
var ErrFailedToMatchHtmlTemplate = errors.New("TemplateName was set, but it failed to match within HtmlTemplate")

func render(res *Body, mediaType string, w io.Writer, cfg Config) error {
	switch mediaType {
	case SupportedMediaTypeHtml:
		dev.Log("Rendering html...")
		if res.TemplateName != "" && cfg.HtmlTemplate != nil {
			dev.Log("Template name is set, so expecting a template...")

			if tm := cfg.HtmlTemplate.Lookup(res.TemplateName); tm != nil {
				dev.Log("Executing HtmlTemplate...")
				err := tm.ExecuteTemplate(w, res.TemplateName, res.Data)
				if err != nil {
					return fmt.Errorf("rendering data in html template %s: %w", res.TemplateName, err)
				}
				break
			}

			return ErrFailedToMatchHtmlTemplate
		}
		return fmt.Errorf("not using a template because either HtmlTemplate or TemplateName is not set")
	case SupportedMediaTypePlaintext:
		dev.Log("Rendering plain text...")

		if res.TemplateName != "" && cfg.TextTemplate != nil {
			dev.Log("Template name is set, so expecting a template...")

			if tm := cfg.TextTemplate.Lookup(res.TemplateName); tm != nil {
				dev.Log("Executing TextTemplate...")
				err := tm.ExecuteTemplate(w, res.TemplateName, res.Data)
				if err != nil {
					return fmt.Errorf("rendering data in text template %s: %w", res.TemplateName, err)
				}
				break
			}

			return ErrFailedToMatchTextTemplate
		}
		dev.Log("Not using a template because either TextTemplate or TemplateName is not set...")

		if data, ok := res.Data.(string); ok {
			dev.Log("Can write data directly because it is a string...")
			_, err := w.Write([]byte(data))
			if err != nil {
				return fmt.Errorf("rendering data as plain string: %w", err)
			}
			break
		}

		return fmt.Errorf("trying to render data as %v but this type is not supported: %#v", SupportedMediaTypePlaintext, res.Data)
	case SupportedMediaTypeJson:
		dev.Log("Rendering json...")
		enc := json.NewEncoder(w)
		enc.SetIndent(cfg.JsonPrefix, cfg.JsonIndent)
		err := enc.Encode(res.Data)
		if err != nil {
			return fmt.Errorf("rendering data as JSON: %w", err)
		}
	case SupportedMediaTypeXml:
		dev.Log("Rendering xml...")
		enc := xml.NewEncoder(w)
		enc.Indent(cfg.XmlPrefix, cfg.XmlIndent)
		err := enc.Encode(res.Data)
		if err != nil {
			return fmt.Errorf("rendering data as XML: %w", err)
		}
	case SupportedMediaTypeCsv:
		data, ok := res.Data.(Csv)
		if !ok {
			return fmt.Errorf("trying to write %#v, but it does not implement rsvp.Csv", res.Data)
		}
		dev.Log("Rendering csv...")
		wr := csv.NewWriter(w)
		err := data.MarshalCsv(wr)
		if err != nil {
			return fmt.Errorf("rendering data as CSV: %w", err)
		}
	case SupportedMediaTypeBytes:
		dev.Log("Rendering bytes...")
		_, err := w.Write(res.Data.([]byte))
		if err != nil {
			return fmt.Errorf("rendering data as bytes: %w", err)
		}
	case SupportedMediaTypeGob:
		dev.Log("Rendering gob...")
		err := gob.NewEncoder(w).Encode(res.Data)
		if err != nil {
			return fmt.Errorf("rendering data as encoding/gob: %w", err)
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
