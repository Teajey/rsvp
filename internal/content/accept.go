package content

import (
	"iter"
	"slices"
	"strings"
)

func splitAndTrimSpace(s, sep string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, split := range strings.Split(s, sep) {
			if !yield(strings.TrimSpace(split)) {
				return
			}
		}
	}
}

func parseAcceptUnsorted(accept string) iter.Seq[proposal] {
	return func(yield func(proposal) bool) {
		for proposal := range splitAndTrimSpace(accept, ",") {
			parsed, err := ParseProposal(proposal)
			if err != nil {
				continue
			}
			if !yield(parsed) {
				return
			}
		}
	}
}

// Returns mediatypes in the order specified by https://httpwg.org/specs/rfc9110.html#field.accept starting with highest precedence.
//
// An empty string will yield a single "*/*".
//
// Invalid proposals are dropped.
func ParseAccept(accept string) iter.Seq[string] {
	if accept == "" {
		return func(yield func(string) bool) {
			yield("*/*")
		}
	}

	list := slices.Collect(parseAcceptUnsorted(accept))
	slices.SortFunc(list, proposalCmp)

	return func(yield func(string) bool) {
		for _, p := range list {
			if !yield(p.MediaType()) {
				return
			}
		}
	}
}
