# RSVP

My "functional" wrapper around Golang's [`net/http` server stuff](https://pkg.go.dev/net/http#hdr-Servers)

It's easy for me to lose track of what I've written to [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter) occasionally receiving the old
```
http: multiple response.WriteHeader calls
```
With this library I just return a value to trigger the HTTP response write.

## Features
 - Respects the Accept header and will attempt to provide the data in a format that is requested.
   - [x] JSON (default)
   - [x] HTML
   - [ ] Plain text
   - [ ] Other?
