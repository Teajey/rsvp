package content

import (
	"iter"
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

func ParseAccept(accept string) {
	// proposals := make([]Proposal, 0)
	// for proposal := range splitAndTrimSpace(accept, ",") {
	// 	proposals = append(proposals, ParseProposal())
	// }
}
