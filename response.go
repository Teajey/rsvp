package rsvp

import (
	"encoding/json"
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
	// Beware that the default value of nil will likely render as JSON "null\n" rather
	// than the expected empty body. Set Body to "" to avoid this.
	Body         any
	TemplateName string
	Status       int

	predeterminedMediaType   supportedType
	predeterminedContentType string

	seeOther          string
	movedPermanently  string
	permanentRedirect string
}

func (res *Response) mediaTypes(cfg *Config) iter.Seq[supportedType] {
	return func(yield func(supportedType) bool) {
		if res.predeterminedMediaType != "" {
			log.Dev("Overriding media-types with %s", res.predeterminedMediaType)
			yield(supportedType(res.predeterminedMediaType))
			return
		}

		switch res.Body.(type) {
		case string:
			if !yield(mPlaintext) {
				return
			}
		case []byte:
			yield(mBytes)
			return
		}

		if !yield(mJson) {
			return
		}

		if res.TemplateName != "" {
			if cfg.HtmlTemplate != nil && cfg.HtmlTemplate.Lookup(res.TemplateName) != nil {
				if !yield(mHtml) {
					return
				}
			}

			if cfg.TextTemplate != nil && cfg.TextTemplate.Lookup(res.TemplateName) != nil {
				if !yield(mPlaintext) {
					return
				}
			}
		}
	}
}

type Config struct {
	HtmlTemplate *html.Template
	TextTemplate *text.Template

	// Controls which file extensions override the Accept header. E.g. "json" will only accept "application/json" by default.
	//
	// You might instead set "json" to accept "application/*", or "*/*" (although the latter is the default if "json" weren't set at all)
	ExtToProposalMap map[string]string
}

// Sets Config.ExtensionToProposalMap = defaultExtToProposalMap
func DefaultConfig() *Config {
	return &Config{
		ExtToProposalMap: defaultExtToProposalMap,
	}
}

func (res *Response) Write(w http.ResponseWriter, r *http.Request, cfg *Config) error {
	if res.movedPermanently != "" {
		http.Redirect(w, r, res.movedPermanently, http.StatusMovedPermanently)
		return nil
	}

	if res.permanentRedirect != "" {
		http.Redirect(w, r, res.permanentRedirect, http.StatusPermanentRedirect)
		return nil
	}

	h := w.Header()

	accept := r.Header.Get("Accept")

	supported := slices.Collect(res.mediaTypes(cfg))
	log.Dev("supported %v", supported)

	var ext string
	if r.Method == http.MethodGet {
		ext = strings.TrimPrefix(filepath.Ext(r.URL.Path), ".")
	}

	mediaType := chooseMediaType(ext, supported, content.ParseAccept(accept), cfg.ExtToProposalMap)
	log.Dev("mediaType %#v", mediaType)

	if mediaType == "" {
		if ext != "" {
			w.WriteHeader(http.StatusNotFound)
			return nil
		}

		w.WriteHeader(http.StatusNotAcceptable)
		return nil
	}

	if mediaType == "text/plain" && cfg.TextTemplate == nil && res.TemplateName != "" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	if res.Body == nil {
		log.Dev("Early returning because body is nil")
		return nil
	}

	contentType := h.Get("Content-Type")
	if contentType == "" {

		if res.predeterminedMediaType != "" {
			// In this case, assuming mediaType == res.predeterminedMediaType
			contentType = res.predeterminedContentType
		} else {
			contentType = mediaTypeToContentType[mediaType]
		}

		log.Dev("Setting content-type to %#v", contentType)
		h.Set("Content-Type", contentType)
	}

	if res.Status != 0 {
		w.WriteHeader(res.Status)
	}

	// If the client's getting HTML they're probably using a browser which will
	// automatically follow this SeeOther. We shouldn't bother rendering anything.
	if res.seeOther != "" && mediaType == "text/html" {
		http.Redirect(w, r, res.seeOther, http.StatusSeeOther)
		return nil
	}

	switch mediaType {
	case mHtml:
		log.Dev("Rendering html...")
		if res.TemplateName != "" && cfg.HtmlTemplate != nil {
			err := cfg.HtmlTemplate.ExecuteTemplate(w, res.TemplateName, res.Body)
			if err != nil {
				return fmt.Errorf("Failed to render body HTML template %s: %w", res.TemplateName, err)
			}
		} else {
			_, err := w.Write([]byte(res.Body.(string)))
			if err != nil {
				return fmt.Errorf("Failed to write string as HTML: %w", err)
			}
		}
	case mPlaintext:
		log.Dev("Rendering plain text...")

		if res.TemplateName != "" {
			log.Dev("Template name is set, so expecting a template...")

			if tm := cfg.TextTemplate.Lookup(res.TemplateName); tm != nil {
				log.Dev("Executing TextTemplate...")
				err := tm.ExecuteTemplate(w, res.TemplateName, res.Body)
				if err != nil {
					return fmt.Errorf("Failed to render body as text template %s: %w", res.TemplateName, err)
				}
				break
			}

			return fmt.Errorf("TemplateName was set, but there is no TextTemplate to check")
		}
		log.Dev("Not using a template because either TextTemplate or TemplateName is not set...")

		if body, ok := res.Body.(string); ok {
			log.Dev("Can write text directly because it is a string...")
			_, err := w.Write([]byte(body))
			if err != nil {
				return fmt.Errorf("Failed to render body as plain-text string: %w", err)
			}
			break
		}

		return fmt.Errorf("Trying to render body as %s but this type is not supported for strings: %#v", mPlaintext, res.Body)
	case mJson:
		log.Dev("Rendering json...")
		err := json.NewEncoder(w).Encode(res.Body)
		if err != nil {
			return fmt.Errorf("Failed to render body as JSON: %w", err)
		}
	case mBytes:
		log.Dev("Rendering bytes...")
		_, err := w.Write(res.Body.([]byte))
		if err != nil {
			return fmt.Errorf("Failed to render body as bytes: %w", err)
		}
	default:
		return fmt.Errorf("Unhandled mediaType: %#v", mediaType)
	}

	if res.seeOther != "" {
		http.Redirect(w, r, res.seeOther, http.StatusSeeOther)
		return nil
	}

	return nil
}

// Will redirect to the given URL after writing the response body.
//
// Except for rendering HTML, the body will still be rendered in
// this case. For instance, if the request was a JSON PUT from
// the commandline it's helpful to see the result without having
// to manually follow the Location header.
func SeeOther(url string) Response {
	return Response{
		seeOther: url,

		// Status must be set here otherwise any writes
		// to http.ResponseWriter will beat us to the punch
		Status: http.StatusSeeOther,
	}
}

// Will perform an immediate 301 using the given URL.
func MovedPermanently(url string) Response {
	return Response{movedPermanently: url}
}

// Will perform an immediate 308 using the given URL.
//
// 308 is intended for non-GET links/operations.
//
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Redirections#permanent_redirections
func PermanentRedirect(url string) Response {
	return Response{permanentRedirect: url}
}

// Short-hand for returning rsvp.Response{} which is equivalent to a blank 200 OK response
func Ok() Response {
	return Response{}
}

// Set body to html using a string, making sure "Content-Type: text/html; charset=utf-8" is set.
//
// Use rsvp.ServeMux.HtmlTemplate and rsvp.Response.TemplateName to render from an HTML template.
func (r *Response) Html(html string) {
	r.Body = html
	r.predeterminedMediaType = mHtml
	r.predeterminedContentType = mediaTypeToContentType[mHtml]
}
