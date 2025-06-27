package rsvp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Body     any
	Template string
	SeeOther string
	Status   int
}

func (res *Response) Write(w http.ResponseWriter, r *http.Request, t *template.Template) error {
	if res.SeeOther != "" {
		http.Redirect(w, r, res.SeeOther, http.StatusSeeOther)
		return nil
	}

	bodyBytes := bytes.NewBuffer([]byte{})
	accept := r.Header.Get("Accept")

	if res.Status != 0 {
		w.WriteHeader(res.Status)
	}

	// I'm too dumb to get my head around the accept header's weighting feature. So I just pick the first match, for now
	switch {
	case strings.Contains(accept, "text/html"):
		if t == nil {
			err := json.NewEncoder(bodyBytes).Encode(res.Body)
			if err != nil {
				return err
			}
		} else {
			subTemplate := t.Lookup(res.Template)
			if subTemplate == nil {
				err := t.Execute(bodyBytes, res.Body)
				if err != nil {
					return err
				}
			} else {
				err := subTemplate.Execute(bodyBytes, res.Body)
				if err != nil {
					return err
				}
			}
		}
	case strings.Contains(accept, "application/json"):
		err := json.NewEncoder(bodyBytes).Encode(res.Body)
		if err != nil {
			return err
		}
	default:
		switch data := res.Body.(type) {
		case io.Reader:
			_, err := bodyBytes.ReadFrom(data)
			if err != nil {
				return err
			}
		default:
			return errors.New("Unsupported response body type")
		}
	}

	_, err := w.Write(bodyBytes.Bytes())
	if err != nil {
		return fmt.Errorf("Failed to write rsvp.Response bodyBytes to http.ResponseWriter: %s", err)
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
		res.Template = template[0]
	}
	return res
}

func SeeOther(url string) Response {
	return Response{SeeOther: url}
}
