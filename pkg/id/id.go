package id

import (
	"regexp"
	"strings"
)

var bracket, invalid, repeat *regexp.Regexp

// Init the regexp
func Init() {
	bracket = regexp.MustCompile(`(\(.*\))|(\[.*\])|(\{.*\})|(<.*>)`)
	invalid = regexp.MustCompile(`[^a-zA-Z0-9-]`)
	repeat = regexp.MustCompile(`-+`)
}

// Identify must call after Init()
func Identify(s string) string {
	s = string(bracket.ReplaceAll([]byte(s), []byte("")))
	s = string(invalid.ReplaceAll([]byte(s), []byte("-")))
	s = string(repeat.ReplaceAll([]byte(s), []byte("-")))

	if len(s) == 0 {
		return "section"
	}

	s = strings.ToLower(s)
	if s[0] < 'a' || s[0] > 'z' {
		s = "section-" + s
	}

	l := len(s)
	if l >= 1 && s[l-1] == '-' {
		if l == 1 {
			return "section"
		}
		s = s[:l-1]
	}
	return s
}
