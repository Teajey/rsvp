//go:build !rsvp_logs

package log

func Dev(format string, v ...any) {
	// NO-OP
}
