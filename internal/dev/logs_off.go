//go:build !rsvp_logs

package dev

func Log(format string, v ...any) {
	// NO-OP
}
