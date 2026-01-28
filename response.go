// Package rsvp is a Go web framework built around content negotiation.
//
// The framework automatically negotiates response format based on the Accept
// header, supporting JSON, XML, HTML, plain text, binary, Gob,
// and MessagePack (using -tags=rsvp_msgpack).
// This content negotiation extends to ALL responses, including redirects,
// allowing you to provide rich feedback in many contexts.
//
// This makes rsvp particularly well-suited for APIs that serve multiple clients
// (browsers, mobile apps, CLI tools) and for taking advantage of principles such
// as REST and progressive enhancement.
package rsvp

type Response struct {
	// Data is the raw data of the response payload to be rendered.
	//
	// IMPORTANT: A nil Data renders as JSON "null\n", not an empty response.
	// Use Data: "" for a blank response body.
	Data any
	// TemplateName sets the template that this Response may attempt to select from
	// [Config.HtmlTemplate] or [Config.TextTemplate],
	//
	// [ResponseWriter.DefaultTemplateName] may also be used to set a default once on a handler.
	//
	// It is not an error if a template is not found for one of the two templates; other formats will be attempted.
	//
	// TODO: Perhaps a warning should be issued to stderr if this fails to match on both templates?
	TemplateName string

	statusCode int

	predeterminedMediaType string

	blankBodyOverride bool

	redirectLocation string
}

func (res *Response) isBlank() bool {
	return res.Data == nil && res.blankBodyOverride
}

// Blank will render as a blank response with no Content-Type.
//
// Status 200 by default.
func Blank() Response {
	return Response{blankBodyOverride: true}
}

// Data is a convenience function equivalent to instantiating Response{Data: data}
func Data(data any) Response {
	return Response{Data: data}
}

// Html can be used to set [Response.Data].
//
// The wrapped string will be treated as text/html
// instead of text/plain.
type Html string
