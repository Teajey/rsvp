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
		"text/html",
		"application/json",
		"text/plain",
		"application/xml",
		"text/*",
		"*/*",
	}
	assert.Eq(t, "Expected length", 6, len(mediaTypes))
	log.Printf("mediaTypes: %v", mediaTypes)
	for i := range mediaTypes {
		assert.Eq(t, fmt.Sprintf("Index %d", i), expected[i], mediaTypes[i])
	}
}
