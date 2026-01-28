package rsvp

import (
	html "html/template"
	text "text/template"
)

// Settings for writing the rsvp.Response
type Config struct {
	// [Response.Data] may be passed to this template as data if content negotiation resolves to text/html and [Response.TemplateName] matches via [html.Template.Lookup].
	//
	// If both HtmlTemplate and TextTemplate match [Response.TemplateName], HtmlTemplate takes precedence.
	HtmlTemplate *html.Template
	// [Response.Data] may be passed to this template as data if content negotiation resolves to text/plain and [Response.TemplateName] matches via [text.Template.Lookup].
	//
	// If both HtmlTemplate and TextTemplate match [Response.TemplateName], HtmlTemplate takes precedence.
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
