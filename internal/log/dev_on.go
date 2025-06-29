//go:build dev_logging

package log

import std "log"

func Dev(format string, v ...any) {
	std.Printf(format, v...)
}
