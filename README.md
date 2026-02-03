# rsvp

net/http middleware providing a Handler interface with type-driven content negotiation.

---

The default net/http handler interface:

```go
ServeHTTP(http.ResponseWriter, *http.Request)
```

rsvp's handler interface looks like this:

```go
type Body struct {
    Data any,
}

ServeHTTP(ResponseWriter, *http.Request) Body
```

## Features

- Content Negotiation. rsvp will attempt to provide Data in a supported media type that is requested via the Accept header; or even the URL's file extension in the case of GET requests:
  - [x] `application/json`
  - [x] `text/html`
  - [x] `text/plain`
  - [x] `text/csv` (by implementing the rsvp.Csv interface)
  - [x] `application/octet-stream`
  - [x] `application/xml`
  - [x] `application/vnd.golang.gob` (Golang's [encoding/gob](https://go.dev/blog/gob))
  - [x] `application/vnd.msgpack` (optional extension behind -tags=rsvp_msgpack)
  - [ ] Others to be implemented?
- Extension matching on GET requests:
  - `/users/123` → Returns default media type (determined by the value of Body)
  - `/users/123.json` → Forces `application/json`
  - `/users/123.xml` → Forces `application/xml`
  - `/users/123.csv` → Forces `text/csv`
  - NOTE: Currently, this behaviour is hidden behind net/http's strict path matching. The above examples would require something like `mux.Handle("/users/{filename}")` with middleware that strips out the file extension, and matches the remaining file stem with its respective handler. rsvp does not currently provide a utility for this.

It's easy for me to lose track of what I've written to [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter). Occasionally receiving the old `http: multiple response.WriteHeader calls`

With this library I just return a value, which I can only ever do once, to execute an HTTP response write. Why write responses with a weird mutable reference from goodness knows where? YEUCH!

Having to remember to return separately from resolving the response? \*wretch\*

```go
if r.Method != http.MethodPut {
	http.Error(w, "Use PUT", http.StatusMethodNotAllowed)
	return
}
```

Not with rsvp 🫠

```go
if r.Method != http.MethodPut {
	return rsvp.Data("Use PUT").StatusMethodNotAllowed()
}
```

(Wrapping this with your own convenience method, i.e. `func ErrorMethodNotAllowed(message string) rsvp.Body` is encouraged. You can decide for yourself how errors are represented)

## Comparison

| Feature             | net/http              | Gin / Echo / Fiber     | rsvp                    |
| ------------------- | --------------------- | ---------------------- | ----------------------- |
| Response Style      | Imperative (w.Write)  | Context-based (c.JSON) | Value-Oriented (return) |
| Content Negotiation | Manual                | Manual / Middleware    | Built-in & Automatic    |
| URL Extensions      | Manual parsing        | Generally unsupported  | Native (.json, .csv)    |
| Response Handling   | Easy to forget return | Side-effect based      | Compile-time enforced   |
| Size                | Standard minimum      | Massive abstraction    | Lightweight Wrapper     |

## When to use rsvp

You value Progressive Enhancement: You want your API to be easily browsable by a human (HTML/XML) but consumable by a script (JSON/CSV) using the same URL.

You hate "Multiple WriteHeader" logs: You want a handler signature that makes it impossible to write a partial or double response.

You want trivially testable handlers: Since handlers return a struct, you can unit test your logic by inspecting the returned rsvp.Response instead of mocking a whole http.ResponseWriter.

## When to stick with net/http or others

High-performance binary streaming: If you are streaming gigabytes of data where every nanosecond of overhead matters, the abstraction of rsvp struct might not be for you.

OpenAPI-first workflows: If your primary goal is generating documentation from code, a schema-heavy framework like Huma might be a better fit. Although you might want to look into creating a proper RESTful self-documenting interface by combining rsvp with hyprctl: https://github.com/Teajey/hyprctl

## Quickstart

```go
func main() {
    mux := http.NewServeMux()
    cfg := rsvp.Config{}
    mux.HandleFunc("GET /users/{id}", rsvp.AdaptHandlerFunc(cfg, getUser))
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
rsvp.Config{
    HtmlTemplate: template.Must(template.ParseGlob("templates/html/*.gotmpl")),
    TextTemplate: template.Must(template.ParseGlob("templates/text/*.gotmpl")),
}

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
