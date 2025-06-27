# RSVP

My "functional" wrapper around Golang's [`net/http` server stuff](https://pkg.go.dev/net/http#hdr-Servers)

It's easy for me to lose track of what I've written to [`http.ResponseWriter`](https://pkg.go.dev/net/http#ResponseWriter) occasionally receiving the old
```
http: multiple response.WriteHeader calls
```
Instead, I just to return a value to trigger the HTTP response write.
