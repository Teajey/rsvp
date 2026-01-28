// Package rsvp is a Go web framework built around content negotiation.
//
// The framework automatically negotiates response format based on the Accept
// header, supporting JSON, XML, HTML, plain text, binary, Gob,
// and MessagePack (using -tags=rsvp_msgpack).
// This content negotiation extends to ALL responses, including redirects,
// allowing you to provide rich feedback in many contexts.
//
// The Accept header should be expected to be used as standardized; weighting is supported.
// If an acceptable fallback is not reached, a 406 Not Acceptable will be returned, and
// the Content-Type and body will be set as if Accept: */* was sent.
//
// This makes rsvp particularly well-suited for APIs that serve multiple clients
// (browsers, mobile apps, CLI tools) and for taking advantage of principles such
// as REST and progressive enhancement.
package rsvp

// Body represents the content body of an HTTP response.
//
// By default, it represents a 200 OK response. The Body.Status* methods (e.g. [Body.StatusFound]) may be used to set a non-200 status.
type Body struct {
	// Data is the raw data of the response payload to be rendered.
	//
	// IMPORTANT: A nil Data renders as JSON "null\n", not an empty response.
	// Use Data("") for a blank text/plain response body, or [Blank] for a blank response with no Content-Type.
	Data any
	// TemplateName sets the template that this Body may attempt to select from
	// [Config.HtmlTemplate] or [Config.TextTemplate],
	//
	// [ResponseWriter.DefaultTemplateName] may also be used to set a default once on a handler.
	//
	// It is not an error if a template is not found for one of the two templates; other formats will be attempted.
	TemplateName string
	// TODO: Perhaps a warning should be issued to stderr if this fails to match on both templates?

	statusCode int

	predeterminedMediaType string

	blankBodyOverride bool

	redirectLocation string
}

func (res *Body) isBlank() bool {
	return res.Data == nil && res.blankBodyOverride
}

// Blank will render as a blank response with no Content-Type.
//
// Status 200 by default.
func Blank() Body {
	return Body{blankBodyOverride: true}
}

// Data is a convenience function equivalent to instantiating Body{Data: data}
//
// IMPORTANT: nil Data renders as JSON "null\n", not an empty response.
// Use Data("") for a blank text/plain response body, or [Blank] for a blank response with no Content-Type.
func Data(data any) Body {
	return Body{Data: data}
}
