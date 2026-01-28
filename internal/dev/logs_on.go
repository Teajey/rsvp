//go:build rsvp_logs

package dev

import "log"

func Log(format string, v ...any) {
	log.Printf(format, v...)
}
