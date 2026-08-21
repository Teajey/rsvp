package rsvp

import (
	"bytes"
	"errors"
	"testing"

	"github.com/Teajey/rsvp/internal/assert"
)

// countingFlusher is a minimal http.Flusher fake that just counts calls.
type countingFlusher struct {
	flushes int
}

func (f *countingFlusher) Flush() {
	f.flushes++
}

// erroringWriter always fails, to test that flushWriter propagates
// the underlying writer's error untouched.
type erroringWriter struct {
	err error
}

func (w erroringWriter) Write(p []byte) (int, error) {
	return 0, w.err
}

// A successful write that reaches the threshold should trigger exactly
// one Flush.
func TestFlushWriterFlushesOnSuccessfulWriteAtThreshold(t *testing.T) {
	var buf bytes.Buffer
	flusher := &countingFlusher{}
	fw := &flushWriter{w: &buf, f: flusher, threshold: 5}

	n, err := fw.Write([]byte("hello"))
	assert.FatalErr(t, "Write", err)
	assert.Eq(t, "bytes written", 5, n)

	assert.Eq(t, "buffer contents", "hello", buf.String())
	assert.Eq(t, "flush count", 1, flusher.flushes)
	assert.Eq(t, "buffered resets after flush", 0, fw.buffered)
}

func TestFlushWriterNoFlushBelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	flusher := &countingFlusher{}
	fw := &flushWriter{w: &buf, f: flusher, threshold: 100}

	_, err := fw.Write([]byte("hi"))
	assert.FatalErr(t, "Write #1", err)
	_, err = fw.Write([]byte(" there"))
	assert.FatalErr(t, "Write #2", err)

	assert.Eq(t, "flush count", 0, flusher.flushes)
	assert.Eq(t, "buffered", 8, fw.buffered)
}

// StreamEager() leaves StreamingThreshold at its zero value, so every
// successful write should flush immediately.
func TestFlushWriterZeroThresholdFlushesEveryWrite(t *testing.T) {
	var buf bytes.Buffer
	flusher := &countingFlusher{}
	fw := &flushWriter{w: &buf, f: flusher} // threshold defaults to 0

	for i, chunk := range []string{"a", "b", "c"} {
		_, err := fw.Write([]byte(chunk))
		assert.FatalErr(t, "Write", err)
		assert.Eq(t, "flush count after write", i+1, flusher.flushes)
	}
}

func TestFlushWriterAccumulatesThenFlushesAtThreshold(t *testing.T) {
	var buf bytes.Buffer
	flusher := &countingFlusher{}
	fw := &flushWriter{w: &buf, f: flusher, threshold: 10}

	_, err := fw.Write([]byte("abcd")) // 4 buffered, below threshold
	assert.FatalErr(t, "Write #1", err)
	assert.Eq(t, "flush count after first write", 0, flusher.flushes)
	assert.Eq(t, "buffered after first write", 4, fw.buffered)

	_, err = fw.Write([]byte("efghij")) // +6 = 10, hits threshold
	assert.FatalErr(t, "Write #2", err)
	assert.Eq(t, "flush count after second write", 1, flusher.flushes)
	assert.Eq(t, "buffered resets after flush", 0, fw.buffered)
}

func TestFlushWriterPropagatesUnderlyingWriteError(t *testing.T) {
	wantErr := errors.New("boom")
	flusher := &countingFlusher{}
	fw := &flushWriter{w: erroringWriter{err: wantErr}, f: flusher, threshold: 1}

	n, err := fw.Write([]byte("x"))
	assert.Eq(t, "bytes written", 0, n)
	assert.FatalErrIs(t, "Write", err, wantErr)
}
