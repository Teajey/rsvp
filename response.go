package rsvp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	HtmlTemplate *template.Template
	Data         any
	SeeOther     string
	Session      *Session
	Status       int
	LogMessage   string
}

func (res *Response) Write(w http.ResponseWriter, r *http.Request) error {
	if res.LogMessage != "" {
		log.Printf("%s --- %s\n", res.Data, res.LogMessage)
	}

	if res.Session != nil {
		err := res.Session.inner.Save(r, w)
		if err != nil {
			return fmt.Errorf("Failed to write session: %w", err)
		}
	}

	if res.SeeOther != "" {
		http.Redirect(w, r, res.SeeOther, http.StatusSeeOther)
		if res.Session != nil {
			flash := res.Session.FlashPeek()
			if flash != nil {
				err := json.NewEncoder(w).Encode(flash)
				if err != nil {
					return err
				}
			}
		}
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

func (r Response) Log(format string, a ...any) Response {
	r.LogMessage = fmt.Sprintf(format, a...)
	return r
}

func (r Response) LogError(err error) Response {
	r.LogMessage = fmt.Sprintf("%s", err)
	return r
}

func (r Response) SaveSession(s Session) Response {
	r.Session = &s
	return r
}

func SeeOther(url string) Response {
	return Response{SeeOther: url}
}
