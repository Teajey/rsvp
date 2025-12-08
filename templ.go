//go:build rsvp_templ

package rsvp

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

var _ Response = (*Templ)(nil)

type Templ struct {
	Rsvp
	htmlRenderer func(body any) templ.Component
}

func (res Rsvp) Templ(renderer func(body any) templ.Component) Templ {
	return Templ{
		res,
		renderer,
	}
}

func (res *Templ) renderHtml(w http.ResponseWriter, r *http.Request, cfg *Config) error {
	err := res.htmlRenderer(res.Rsvp.Body).Render(r.Context(), w)
	if err != nil {
		return fmt.Errorf("failed to render HTML Templ template %s: %w", res.TemplateName, err)
	}
	return nil
}
