// Package fyslices provides helper functions for working with []string
// that are not available in the standard slices package.
package fyslices

import (
	"fmt"
	"strconv"
	"strings"
)

// From returns s[i:] or an empty slice if `i` is out of range.
// Negative `i` is treated as 0; `i` beyond length is treated as len(s).
func From(s []string, i int) []string {
	if i < 0 {
		i = 0
	}
	if i > len(s) {
		i = len(s)
	}
	return s[i:]
}

// Max returns s[:to] or the whole slice if to is out of range.
// Negative `to` is treated as 0; `to` beyond length is treated as len(s).
func Max(s []string, to int) []string {
	if to < 0 {
		to = 0
	}
	if to > len(s) {
		to = len(s)
	}
	return s[:to]
}

// UpTo returns s[:i+1], or the whole slice if `i` >= len(s).
// Negative `i` returns empty slice.
func UpTo(s []string, i int) []string {
	if i < 0 {
		return []string{}
	}
	if i >= len(s) {
		return s
	}
	return s[:i+1]
}

// Map applies fn to every element and returns a new slice.
func Map(s []string, fn func(string) string) []string {
	out := make([]string, len(s))
	for i, v := range s {
		out[i] = fn(v)
	}
	return out
}

// JoinQuoted joins elements with sep and wraps each element in double quotes.
// Example: []string{"a","b"} with sep="," -> `"a","b"`.
func JoinQuoted(s []string, sep string) string {
	out := make([]string, len(s))
	for i, v := range s {
		out[i] = fmt.Sprintf("%q", v)
	}
	return strings.Join(out, sep)
}

// Filter returns a new slice containing only elements that satisfy fn.
func Filter(s []string, fn func(string) bool) []string {
	out := make([]string, 0, len(s))
	for _, v := range s {
		if fn(v) {
			out = append(out, v)
		}
	}
	return out
}

// Unique removes duplicate values, preserving the first occurrence order.
func Unique(s []string) []string {
	seen := make(map[string]struct{})
	out := []string{}
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

// ContainsAll reports whether at least one of vals is present in s.
func ContainsAll(s []string, vals ...string) bool {
	set := make(map[string]struct{}, len(s))
	for _, v := range s {
		set[v] = struct{}{}
	}
	for _, v := range vals {
		if _, ok := set[v]; ok {
			return true
		}
	}
	return false
}

// Types classifies each element as "int", "bool", or "string".
// Evaluation order: bool first (narrower), then int, finally string.
func Types(s []string) []string {
	res := make([]string, len(s))
	for i, v := range s {
		if _, err := strconv.ParseBool(v); err == nil {
			res[i] = "bool"
		} else if _, err := strconv.Atoi(v); err == nil {
			res[i] = "int"
		} else {
			res[i] = "string"
		}
	}
	return res
}

// ToInt parses s[i] as int.
func ToInt(s []string, i int) (int, error) {
	if i < 0 || i >= len(s) {
		return 0, fmt.Errorf("index %d out of range", i)
	}
	return strconv.Atoi(s[i])
}

// ToBool parses s[i] as bool.
func ToBool(s []string, i int) (bool, error) {
	if i < 0 || i >= len(s) {
		return false, fmt.Errorf("index %d out of range", i)
	}
	return strconv.ParseBool(s[i])
}

// ToString returns s[i] directly.
func ToString(s []string, i int) (string, error) {
	if i < 0 || i >= len(s) {
		return "", fmt.Errorf("index %d out of range", i)
	}
	return s[i], nil
}

// To is a thin generic wrapper that delegates to the helpers above.
// T must be int, bool, or string; otherwise it panics at runtime.
func To[T ~int | ~bool | ~string](s []string, i int) (T, error) {
	var zero T
	switch ptr := any(&zero).(type) {
	case *int:
		v, err := ToInt(s, i)
		*ptr = v
		return zero, err
	case *bool:
		v, err := ToBool(s, i)
		*ptr = v
		return zero, err
	case *string:
		v, err := ToString(s, i)
		*ptr = v
		return zero, err
	default:
		panic(fmt.Sprintf("unsupported type: %T", zero))
	}
}
