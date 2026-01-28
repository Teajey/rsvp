# rsvp

My value-oriented wrapper around Golang's [`net/http` server stuff](https://pkg.go.dev/net/http#hdr-Servers).

The default `net/http` handler interface:

```go
ServeHTTP(http.ResponseWriter, *http.Request)
```

rsvp's handler interface:

```go
ServeHTTP(rsvp.ResponseWriter, *http.Request) rsvp.Body
```

## Features
 - Content Negotiation. rsvp will attempt to provide the data in a supported media-type that is requested via the Accept header, or even the URL's file extension in the case of GET requests:
   - [x] `application/json`
   - [x] `text/html`
   - [x] `text/plain`
   - [x] `text/csv` (by implementing the rsvp.Csv interface)
   - [x] `application/octet-stream`
   - [x] `application/xml`
   - [x] `application/vnd.golang.gob` (Golang's [encoding/gob](https://go.dev/blog/gob))
   - [x] `application/vnd.msgpack` (optional extension behind -tags=rsvp_msgpack)
   - [ ] Others?
 - Extension matching on GET requests:
   - `/users/123` → Returns default media type (determined by the value of Body)
   - `/users/123.json` → Forces `application/json`
   - `/users/123.xml` → Forces `application/xml`
   - `/users/123.csv` → Forces `text/csv`

It's easy for me to lose track of what I've written to [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter). Occasionally receiving the old `http: multiple response.WriteHeader calls`

With this library I just return a value, which I can only ever do once, to execute an HTTP response write. Why write responses with a weird mutable reference from goodness knows where? YEUCH!

Having to remember to return separately from resolving the response? \*wretch*

```go
if r.Method != http.MethodPut {
	http.Error(w, "Use PUT", http.StatusMethodNotAllowed)
	return
}
```

Not with rsvp 🫠

```go
if r.Method != http.MethodPut {
	return rsvp.Body{Data: "Use PUT"}.StatusMethodNotAllowed()
}
```

(Wrapping this with your own convenience method, i.e. `func ErrorMethodNotAllowed(message string) rsvp.Body` is encouraged. You can decide for yourself how errors are represented)

## Quickstart

```go
func main() {
    mux := rsvp.NewServeMux()
    mux.HandleFunc("GET /users/{id}", getUser)
    http.ListenAndServe(":8080", mux)
}

func getUser(w rsvp.ResponseWriter, r *http.Request) rsvp.Body {
    return rsvp.Data(User{ID: 123}) // In content negotiation this will be offered as, in order; JSON, XML, and encoding/gob.
}
```

> [!IMPORTANT]
> nil Data renders as JSON "null\n", not an empty response.
> 
> Use `rsvp.Data("")` for a blank text/plain response body, or `rsvp.Blank()` for a blank response with no Content-Type.


## Examples

### Templates

```go
mux.Config.HtmlTemplate = template.Must(template.ParseGlob("templates/html/*.gotmpl"))
mux.Config.TextTemplate = template.Must(template.ParseGlob("templates/text/*.gotmpl"))

func showUser(w rsvp.ResponseWriter, r *http.Request) rsvp.Body {
    w.DefaultTemplateName("user.gotmpl") // Must exist in HtmlTemplate and/or TextTemplate for formats to match
    return rsvp.Data(User{ID: 123}) // In content negotiation this will be offered as JSON, XML, HTML, plain text, and encoding/gob.
}
```

### Error responses

Build your own!

```go
type APIError struct {
    Message string `json:"message"`
    Code    string `json:"code"`
}

func ErrorNotFound(msg string) rsvp.Body {
    return rsvp.Data(APIError{Message: msg, Code: "NOT_FOUND"}).StatusNotFound()
}
```

### CSV

```go
type UserList []User

var users UserList

func (ul UserList) MarshalCsv(w *csv.Writer) error {
    w.Write([]string{"ID", "Name", "Email"})
    for _, u := range ul {
        w.Write([]string{u.ID, u.Name, u.Email})
    }
    return nil
}

func userList(w rsvp.ResponseWriter, r *http.Request) rsvp.Body {
    return rsvp.Data(users) // In content negotiation this will be offered as JSON, XML, CSV, and encoding/gob.
}
```

### net/http middleware compatibility

See [middleware_test.go](./middleware_test.go) for an example of how to use this library with standard middleware.

### Live

You can see it in action on my stupid little blog site, brightscroll.net. For instance, https://brightscroll.net/posts/2025-06-30.md vs. https://brightscroll.net/posts/2025-06-30.md.txt
