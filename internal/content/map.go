package content

import (
	"iter"
	"strings"
)

func parseParameters(s string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for p := range splitAndTrimSpace(s, ";") {
			if !yield(parsePair(p)) {
				return
			}
		}
	}
}

func parsePair(p string) (string, string) {
	pairs := strings.SplitN(p, "=", 2)
	key := strings.TrimSpace(pairs[0])
	var value string
	if len(pairs) > 1 {
		value = strings.TrimSpace(pairs[1])
	} else {
		value = ""
	}
	return key, value
}

func insertAppend(m map[string][]string, seq iter.Seq2[string, string]) {
	for k, v := range seq {
		m[k] = append(m[k], v)
	}
}

func collectAppend(seq iter.Seq2[string, string]) map[string][]string {
	m := make(map[string][]string)
	insertAppend(m, seq)
	return m
}
