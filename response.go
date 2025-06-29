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

func (res *Response) MediaTypes(html *html.Template, text *text.Template) iter.Seq[supportedType] {
	return func(yield func(supportedType) bool) {
		if html != nil && html.Lookup(res.TemplateName) != nil {
			if !yield(mHtml) {
				return
			}
		}

		if text != nil && text.Lookup(res.TemplateName) != nil {
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

func (res *Response) Write(w http.ResponseWriter, r *http.Request, h *html.Template, t *text.Template) error {
	if res.SeeOther != "" {
		http.Redirect(w, r, res.SeeOther, http.StatusSeeOther)
		return nil
	}

	accept := r.Header.Get("Accept")

	if res.Status != 0 {
		w.WriteHeader(res.Status)
	}

	supported := slices.Collect(res.MediaTypes(h, t))
	mediaType := resolveMediaType(r.URL, supported, content.ParseAccept(accept))

	switch mediaType {
	case string(mHtml):
		err := h.ExecuteTemplate(w, res.TemplateName, res.Body)
		if err != nil {
			return err
		}
	case string(mPlaintext):
		if t != nil {
			if tm := t.Lookup(res.TemplateName); tm != nil {
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
		w.Header().Set("Content-Type", string(mJson))
		err := json.NewEncoder(w).Encode(res.Body)
		if err != nil {
			return err
		}
	case string(mBytes):
		w.Header().Set("Content-Type", string(mBytes))
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
