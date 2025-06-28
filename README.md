# RSVP

My "functional" wrapper around Golang's [`net/http` server stuff](https://pkg.go.dev/net/http#hdr-Servers)

The default `net/http` interface:

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

RSVP's interface:

```go
type Handler interface {
	ServeHTTP(h http.Header, r *http.Request) Response
}
```

It's easy for me to lose track of what I've written to [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter) occasionally receiving the old `http: multiple response.WriteHeader calls`

With this library I just return a value, and I can only ever do it once, to execute an HTTP response write. Why write responses with a weird mutable reference from goodness knows where? YEUCH!

Having to remember to return separately from resolving the response? \*wretch*

```go
if r.Method != http.MethodPut {
	http.Error(w, "Use PUT", http.StatusNotFound)
	return
}
```

Not with RSVP 🫠

```go
if r.Method != http.MethodPut {
	return rsvp.Response{Status: http.StatusMethodNotAllowed, Body: "Use Put"}
}
```

(Wrapping this with your own convenience method, i.e. `func ErrorMethodNotAllowed(message string) rsvp.Response` is encouraged. You get to decide for yourself how errors are represented)

## Features
 - Respects the Accept header and will attempt to provide the data in a format that is requested.
   - [x] JSON
   - [x] HTML
   - [ ] Plain text
   - [ ] Other?
