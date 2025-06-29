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
	Body         any
	TemplateName string
	SeeOther     string
	Status       int
}

func (res *Response) MediaTypes(cfg *Config) iter.Seq[supportedType] {
	return func(yield func(supportedType) bool) {
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
	h := w.Header()

	if res.SeeOther != "" {
		http.Redirect(w, r, res.SeeOther, http.StatusSeeOther)
		return nil
	}

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

func SeeOther(url string) Response {
	return Response{SeeOther: url}
}

// Short-hand for returning empty rsvp.Response{} which is equivalent to a blank 200 OK response
func Ok() Response {
	return Response{}
}
