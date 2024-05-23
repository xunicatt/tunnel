package logerr

import (
	"fmt"
	"os"
)

const (
	SUCCESS = "\033[32m"
	WARN    = "\033[93m"
	FAIL    = "\033[91m"
	BOLD    = "\033[1m"
	RESET   = "\033[0m"
)

func write(fd *os.File, errType, color string, format *string, args ...any) {
	fmt.Fprintf(fd, color+"tunnel: %s: ", errType)
	fmt.Fprintf(fd, *format+RESET+"\n", args...)
}

func Success(format string, args ...any) {
	write(os.Stdout, "success", SUCCESS, &format, args...)
}

func Message(format string, args ...any) {
	write(os.Stdout, "log", RESET, &format, args...)
}

func Warn(format string, args ...any) {
	write(os.Stderr, "warn", WARN, &format, args...)
}

func Error(format string, args ...any) {
	write(os.Stderr, "error", FAIL, &format, args...)
}
