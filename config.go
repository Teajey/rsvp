package rsvp

import (
	html "html/template"
	text "text/template"
)

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
