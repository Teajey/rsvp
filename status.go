package rsvp

import "net/http"

// Success (2xx)

// StatusCreated sets the response as 201 Created and sets the Location header.
//
// It indicates that a new resource has been successfully created at the given
// location.
func (r Response) StatusCreated(location string) Response {
	r.statusCode = http.StatusCreated
	r.redirectLocation = location
	return r
}

// StatusAccepted sets the response as 202 Accepted.
//
// It indicates that the request has been accepted for processing, but the
// processing has not been completed.
func (r Response) StatusAccepted() Response {
	r.statusCode = http.StatusAccepted
	return r
}

// StatusNoContent sets the response as 204 No Content.
//
// It indicates that the request was successful but there is no content to
// return. Commonly used for DELETE operations or updates with no response body.
func (r Response) StatusNoContent() Response {
	r.statusCode = http.StatusNoContent
	return r
}

// Redirection (3xx)

// StatusMovedPermanently sets the response as 301 Moved Permanently and sets
// the Location header.
//
// Moved Permanently is intended for GET requests.
//
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Redirections#permanent_redirections
func (r Response) StatusMovedPermanently(location string) Response {
	r.statusCode = http.StatusMovedPermanently
	r.redirectLocation = location
	return r
}

// StatusFound sets the response as 302 Found and sets the Location header.
//
// It indicates that the requested resource has been temporarily moved to the
// given location.
func (r Response) StatusFound(location string) Response {
	r.statusCode = http.StatusFound
	r.redirectLocation = location
	return r
}

// StatusSeeOther sets the response as 303 See Other and sets the Location
// header.
//
// See Other is used for redirection in response to POST requests.
func (r Response) StatusSeeOther(location string) Response {
	r.statusCode = http.StatusSeeOther
	r.redirectLocation = location
	return r
}

// StatusNotModified sets the response as 304 Not Modified.
//
// It indicates that the resource has not been modified since the version
// specified by the request headers. Used for conditional requests and caching.
func (r Response) StatusNotModified() Response {
	r.statusCode = http.StatusNotModified
	return r
}

// StatusTemporaryRedirect sets the response as 307 Temporary Redirect and sets
// the Location header.
//
// Temporary Redirect is like 302 Found but guarantees that the HTTP method
// will not be changed when the redirected request is made.
func (r Response) StatusTemporaryRedirect(location string) Response {
	r.statusCode = http.StatusTemporaryRedirect
	r.redirectLocation = location
	return r
}

// StatusPermanentRedirect sets the response as 308 Permanent Redirect and sets
// the Location header.
//
// Permanent Redirect is intended for non-GET requests.
//
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Redirections#permanent_redirections
func (r Response) StatusPermanentRedirect(location string) Response {
	r.statusCode = http.StatusPermanentRedirect
	r.redirectLocation = location
	return r
}

// Client Errors (4xx)

// StatusBadRequest sets the response as 400 Bad Request.
//
// It indicates that the server cannot process the request due to client error,
// such as malformed request syntax or invalid request parameters.
func (r Response) StatusBadRequest() Response {
	r.statusCode = http.StatusBadRequest
	return r
}

// StatusUnauthorized sets the response as 401 Unauthorized.
//
// It indicates that authentication is required and has failed or has not been
// provided.
func (r Response) StatusUnauthorized() Response {
	r.statusCode = http.StatusUnauthorized
	return r
}

// StatusForbidden sets the response as 403 Forbidden.
//
// It indicates that the server understood the request but refuses to authorize
// it. Unlike 401, authenticating will make no difference.
func (r Response) StatusForbidden() Response {
	r.statusCode = http.StatusForbidden
	return r
}

// StatusNotFound sets the response as 404 Not Found.
//
// It indicates that the server cannot find the requested resource.
func (r Response) StatusNotFound() Response {
	r.statusCode = http.StatusNotFound
	return r
}

// StatusMethodNotAllowed sets the response as 405 Method Not Allowed.
//
// It indicates that the request method is known by the server but is not
// supported by the target resource.
func (r Response) StatusMethodNotAllowed() Response {
	r.statusCode = http.StatusMethodNotAllowed
	return r
}

// StatusConflict sets the response as 409 Conflict.
//
// It indicates that the request conflicts with the current state of the
// server, such as attempting to create a duplicate resource.
func (r Response) StatusConflict() Response {
	r.statusCode = http.StatusConflict
	return r
}

// StatusGone sets the response as 410 Gone.
//
// It indicates that the requested resource is no longer available and will not
// be available again. This is a stronger statement than 404 Not Found.
func (r Response) StatusGone() Response {
	r.statusCode = http.StatusGone
	return r
}

// StatusUnprocessableEntity sets the response as 422 Unprocessable Entity.
//
// It indicates that the request was well-formed but was unable to be followed
// due to semantic errors, such as validation failures.
func (r Response) StatusUnprocessableEntity() Response {
	r.statusCode = http.StatusUnprocessableEntity
	return r
}

// StatusTooManyRequests sets the response as 429 Too Many Requests.
//
// It indicates that the user has sent too many requests in a given amount of
// time. Used for rate limiting.
func (r Response) StatusTooManyRequests() Response {
	r.statusCode = http.StatusTooManyRequests
	return r
}

// Server Errors (5xx)

// StatusInternalServerError sets the response as 500 Internal Server Error.
//
// It indicates that the server encountered an unexpected condition that
// prevented it from fulfilling the request.
func (r Response) StatusInternalServerError() Response {
	r.statusCode = http.StatusInternalServerError
	return r
}

// StatusNotImplemented sets the response as 501 Not Implemented.
//
// It indicates that the server does not support the functionality required to
// fulfill the request.
func (r Response) StatusNotImplemented() Response {
	r.statusCode = http.StatusNotImplemented
	return r
}

// StatusServiceUnavailable sets the response as 503 Service Unavailable.
//
// It indicates that the server is currently unable to handle the request due
// to temporary overload or scheduled maintenance.
func (r Response) StatusServiceUnavailable() Response {
	r.statusCode = http.StatusServiceUnavailable
	return r
}
