package rsvp

import (
	"io"
	"net/http"
)

// StreamEager sets r.Streaming = true
func (r Body) StreamEager() Body {
	r.Streaming = true
	return r
}

// StreamEager sets r.Streaming = true and r.StreamingThreshold = threshold
func (r Body) StreamThreshold(threshold int) Body {
	r.Streaming = true
	r.StreamingThreshold = threshold
	return r
}

type flushWriter struct {
	w         io.Writer
	f         http.Flusher
	threshold int
	buffered  int
}

func (fw *flushWriter) Write(p []byte) (int, error) {
	n, err := fw.w.Write(p)
	fw.buffered += n
	if err == nil && fw.buffered >= fw.threshold {
		fw.f.Flush()
		fw.buffered = 0
	}
	return n, err
}
