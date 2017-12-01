package rucksack

import (
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Env returns the value of first not-blank key
func Env(keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return ""
}

// Tags parses tags from a string
func Tags(s string) []string {
	return strings.FieldsFunc(s, func(c rune) bool {
		switch c {
		case ':', '=':
			return false
		default:
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
	})
}

// Fields parses key-value fields from a string
func Fields(s string) map[string]interface{} {
	tags := Tags(s)
	if len(tags) == 0 {
		return nil
	}

	fields := make(map[string]interface{}, len(tags))
	for _, tag := range tags {
		parts := strings.SplitN(tag, ":", 2)
		if len(parts) != 2 || parts[0] == "" {
			continue
		}

		var v interface{} = parts[1]
		if n, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
			v = n
		}
		fields[parts[0]] = v
	}
	return fields
}
