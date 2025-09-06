//go:build rsvp_logs

package log

import std "log"

func Dev(format string, v ...any) {
	std.Printf(format, v...)
}
