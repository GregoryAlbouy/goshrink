package httplog

import "fmt"

var (
	Green = color("\033[1;32m%s\033[0m")
	Red   = color("\033[1;31m%s\033[0m")
)

func color(a string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(a, fmt.Sprint(args...))
	}
	return sprint
}

func statusColor(statusCode int) func(...interface{}) string {
	if statusCode <= 400 {
		return Green
	}
	return Red
}
