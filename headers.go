package rsvp

import "strings"

func contentTypeExtractMediaType(contentType string) string {
	strs := strings.Split(contentType, ";")
	return strings.TrimSpace(strs[0])
}
