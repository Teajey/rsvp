package rsvp

import (
	"encoding/json"
	html "html/template"
	"iter"
	"net/http"
	"slices"
	text "text/template"

	"github.com/Teajey/rsvp/internal/content"
)

type Response struct {
	// Beware that the default value of nil will likely render as JSON "null\n" rather
	// than the expected empty body. Set Body to "" to avoid this.
	Body         any
	TemplateName string
	Status       int

	seeOther          string
	movedPermanently  string
	permanentRedirect string
}

func (res *Response) MediaTypes(cfg *Config) iter.Seq[supportedType] {
	return func(yield func(supportedType) bool) {
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

	if res.Status != 0 {
		w.WriteHeader(res.Status)
	}

	supported := slices.Collect(res.MediaTypes(cfg))
	mediaType := resolveMediaType(r.URL, supported, content.ParseAccept(accept), cfg.ExtToProposalMap)

	switch mediaType {
	case string(mHtml):
		err := cfg.HtmlTemplate.ExecuteTemplate(w, res.TemplateName, res.Body)
		if err != nil {
			return err
		}
	case string(mPlaintext):
		if cfg.TextTemplate != nil {
			if tm := cfg.TextTemplate.Lookup(res.TemplateName); tm != nil {
				err := tm.ExecuteTemplate(w, res.TemplateName, res.Body)
				if err != nil {
					return err
				}
			}
		} else {
			_, err := w.Write([]byte(res.Body.(string)))
			if err != nil {
				return err
			}
		}
	case string(mJson):
		if h.Get("Content-Type") == "" {
			w.Header().Set("Content-Type", string(mJson))
		}
		err := json.NewEncoder(w).Encode(res.Body)
		if err != nil {
			return err
		}
	case string(mBytes):
		if h.Get("Content-Type") == "" {
			w.Header().Set("Content-Type", string(mBytes))
		}
		_, err := w.Write(res.Body.([]byte))
		if err != nil {
			return err
		}
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return nil
	}

	if res.seeOther != "" {
		http.Redirect(w, r, res.seeOther, http.StatusSeeOther)
		return nil
	}

	return nil
}

// Write data as a response body to whatever supported format is requested in the Accept header
// Optionally provide a template name for this response
func Body(data any, template ...string) Response {
	res := Response{
		Body: data,
	}
	if len(template) > 0 {
		res.TemplateName = template[0]
	}
	return res
}

// Will redirect to the given URL after writing the response body.
func SeeOther(url string) Response {
	return Response{seeOther: url}
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

// Short-hand for returning rsvp.Response{Body: ""} which is equivalent to a blank 200 OK response
func Ok() Response {
	return Response{Body: ""}
}

// Short-hand for returning a blank 404 NotFound response.
//
// You can set Body and TemplateName afterwards to add information.
//
// ```
// resp := rsvp.NotFound()
// resp.Body = "404: Couldn't find this page."
// resp.TemplateName = "not_found"
// return resp
// ```
func NotFound() Response {
	return Response{
		Status: http.StatusNotFound,
		Body:   "",
	}
}
