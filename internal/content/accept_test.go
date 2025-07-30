package content_test

import (
	"fmt"
	"log"
	"slices"
	"testing"

	"github.com/Teajey/rsvp/internal/assert"
	"github.com/Teajey/rsvp/internal/content"
)

func TestParseAccept(t *testing.T) {
	mediaTypes := slices.Collect(content.ParseAccept("text/*,application/xml, text/html;q=1, text/plain;format=flowed, */*,application/json;q=0.5"))
	expected := []string{
		"text/plain",
		"text/html",
		"application/xml",
		"text/*",
		"*/*",
		"application/json",
	}
	assert.Eq(t, "Expected length", 6, len(mediaTypes))
	log.Printf("mediaTypes: %v", mediaTypes)
	for i := range mediaTypes {
		assert.Eq(t, fmt.Sprintf("Index %d", i), expected[i], mediaTypes[i])
	}
}

func TestFirefoxAcceptHeader(t *testing.T) {
	sq := content.ParseAccept("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	sl := slices.Collect(sq)
	assert.SlicesEq(t, "firefox accepts", sl, []string{"application/xhtml+xml", "text/html", "application/xml", "*/*"})
}
