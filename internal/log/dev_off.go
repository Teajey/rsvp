//go:build !dev_logging

package log

func Dev(format string, v ...any) {
	// NO-OP
}
