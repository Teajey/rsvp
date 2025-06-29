package rsvp

import (
	"testing"

	"github.com/Teajey/rsvp/internal/assert"
)

func TestMediaTypesEqual(t *testing.T) {
	assert.True(t, "application wildcard match", mediaTypesEqual("application/json", "application/*"))
}
