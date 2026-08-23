package rsvp

import (
	"io"
	"net/http"
	"sync"
	"time"
)

// StreamEager sets r.Streaming = true
//
// See [Body.Streaming] for more info
func (r Body) StreamEager() Body {
	r.Streaming = true
	return r
}

// StreamThreshold sets r.Streaming = true and r.StreamingThreshold = threshold
//
// See [Body.StreamingThreshold] for more info
func (r Body) StreamThreshold(threshold int) Body {
	r.Streaming = true
	r.StreamingThreshold = threshold
	return r
}

func (r Body) StreamInterval(interval time.Duration) Body {
	r.Streaming = true
	r.StreamingInterval = interval
	return r
}

type flushWriter struct {
	mu        sync.Mutex
	w         io.Writer
	f         http.Flusher
	threshold int
	interval  time.Duration
	buffered  int
	timer     *time.Timer
	closed    bool
}

func (fw *flushWriter) Write(p []byte) (int, error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	n, err := fw.w.Write(p)
	fw.buffered += n
	if err != nil {
		return n, err
	}
	if fw.buffered >= fw.threshold {
		fw.flushLocked()
	} else if fw.buffered > 0 {
		fw.armLocked()
	}
	return n, err
}

func (fw *flushWriter) armLocked() {
	if fw.interval <= 0 || fw.closed || fw.timer != nil {
		return
	}
	fw.timer = time.AfterFunc(fw.interval, func() {
		fw.mu.Lock()
		defer fw.mu.Unlock()
		fw.timer = nil
		if fw.buffered > 0 && !fw.closed {
			fw.flushLocked()
		}
	})
}

func (fw *flushWriter) flushLocked() {
	fw.f.Flush()
	fw.buffered = 0
	if fw.timer != nil {
		fw.timer.Stop()
		fw.timer = nil
	}
}

func (fw *flushWriter) close() {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.closed = true
	if fw.timer != nil {
		fw.timer.Stop()
		fw.timer = nil
	}
}
