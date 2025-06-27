package rsvp

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	HtmlTemplate *template.Template
	Data         any
	SeeOther     string
	Status       int
}

func (res *Response) Write(w http.ResponseWriter, r *http.Request) error {
	if res.SeeOther != "" {
		http.Redirect(w, r, res.SeeOther, http.StatusSeeOther)
		return nil
	}

	bodyBytes := bytes.NewBuffer([]byte{})
	accept := r.Header.Get("Accept")

	if res.Status != 0 {
		w.WriteHeader(res.Status)
	}

	switch {
	case strings.Contains(accept, "text/html"):
		if res.HtmlTemplate == nil {
			err := json.NewEncoder(bodyBytes).Encode(res.Data)
			if err != nil {
				return err
			}
		} else {
			err := res.HtmlTemplate.Execute(bodyBytes, res.Data)
			if err != nil {
				return err
			}
		}
	case strings.Contains(accept, "application/json"):
		err := json.NewEncoder(bodyBytes).Encode(res.Data)
		if err != nil {
			return err
		}
	default:
		err := json.NewEncoder(bodyBytes).Encode(res.Data)
		if err != nil {
			return err
		}
	}

	_, err := w.Write(bodyBytes.Bytes())
	if err != nil {
		log.Printf("Failed to write rsvp.Response to HTTP: %s\n", err)
	}
	return nil
}

func Data(htmlTemplate *template.Template, data any) Response {
	return Response{
		HtmlTemplate: htmlTemplate,
		Data:         data,
	}
}

func SeeOther(url string) Response {
	return Response{SeeOther: url}
}
