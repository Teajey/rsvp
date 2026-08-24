package rsvp_test

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"testing"
	text "text/template"
	"time"

	"github.com/Teajey/rsvp"
	"github.com/Teajey/rsvp/internal/assert"
)

// TestStreamingDeliversChunksProgressivelyOverHTTP proves streaming actually
// streams: it asserts the client receives each chunk roughly chunkDelay
// apart, rather than getting the whole body in one burst once the handler
// finishes.
func TestStreamingDeliversChunksProgressivelyOverHTTP(t *testing.T) {
	const chunkDelay = 60 * time.Millisecond
	const numChunks = 3

	cfg := rsvp.Config{}
	cfg.TextTemplate = text.New("chunks").Funcs(text.FuncMap{
		// stand-in for slow/paginated work a real handler might do.
		"sleep": func() string {
			time.Sleep(chunkDelay)
			return ""
		},
	})
	cfg.TextTemplate = text.Must(cfg.TextTemplate.Parse(
		`{{range .}}chunk-{{.}}|{{sleep}}{{end}}`,
	))

	handler := rsvp.NewAdapter(cfg).AdaptFunc(func(w rsvp.ResponseWriter, r *http.Request) rsvp.Body {
		w.DefaultTemplateName("chunks")
		chunks := make([]int, numChunks)
		for i := range chunks {
			chunks[i] = i
		}
		return rsvp.Body{Data: chunks}.StreamEager()
	})

	srv := httptest.NewServer(handler)
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	assert.FatalErr(t, "build request", err)
	req.Header.Set("Accept", "text/plain")

	start := time.Now()
	resp, err := srv.Client().Do(req)
	assert.FatalErr(t, "do request", err)
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	arrivals := make([]time.Duration, 0, numChunks)
	for range numChunks {
		_, err := reader.ReadString('|')
		assert.FatalErr(t, "read chunk", err)
		arrivals = append(arrivals, time.Since(start))
	}

	// If flushing works, chunks trickle in roughly chunkDelay apart.
	// If streaming is broken (e.g. Flush only fires on write error,
	// as currently written), the handler buffers everything and the
	// client only receives all chunks at once, back-to-back, right at
	// the very end -- so this gap would collapse to ~0 for every chunk
	// after the first.
	for i := 1; i < len(arrivals); i++ {
		gap := arrivals[i] - arrivals[i-1]
		assert.True(t, "chunk should arrive with a delay from the previous one", gap >= chunkDelay/2)
	}
}
