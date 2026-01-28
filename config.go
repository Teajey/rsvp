package rsvp

import (
	html "html/template"
	text "text/template"
)

// Settings for writing the rsvp.Body
type Config struct {
	// [Body.Data] may be passed to this template as data if content negotiation resolves to text/html and [Body.TemplateName] matches via [html.Template.Lookup].
	//
	// If both HtmlTemplate and TextTemplate match [Body.TemplateName], HtmlTemplate takes precedence.
	HtmlTemplate *html.Template
	// [Body.Data] may be passed to this template as data if content negotiation resolves to text/plain and [Body.TemplateName] matches via [text.Template.Lookup].
	//
	// If both HtmlTemplate and TextTemplate match [Body.TemplateName], HtmlTemplate takes precedence.
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
