package fyslices

import (
	"fmt"
	"strconv"
)

func From(s []string, i int) []string {
	if i < 0 || i >= len(s) {
		return []string{}
	}
	return s[i:]
}

func Max(s []string, to int) []string {
	if to < 0 || to > len(s) {
		return []string{}
	}
	return s[:to]
}

func UpTo(s []string, i int) []string {
	if i < 0 {
		return []string{}
	}
	if i >= len(s) {
		return s
	}
	return s[:i+1]
}

func IndexOf(s []string, target string) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}

func Map(s []string, fn func(string) string) []string {
	out := make([]string, len(s))
	for i, v := range s {
		out[i] = fn(v)
	}
	return out
}

func Types(args []string) []string {
	types := make([]string, len(args))
	for i, arg := range args {
		if _, err := strconv.Atoi(arg); err == nil {
			types[i] = "int"
		} else if _, err := strconv.ParseBool(arg); err == nil {
			types[i] = "bool"
		} else {
			types[i] = "string"
		}
	}
	return types
}

func To[T any](s []string, i int) (T, error) {
	var zero T
	if i < 0 || i >= len(s) {
		return zero, fmt.Errorf("index out of range")
	}

	val := s[i]
	switch any(zero).(type) {
	case int:
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return zero, err
		}
		return any(parsed).(T), nil
	case bool:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return zero, err
		}
		return any(parsed).(T), nil
	case string:
		return any(val).(T), nil
	default:
		return zero, fmt.Errorf("unsupported type: %T", zero)
	}
}
