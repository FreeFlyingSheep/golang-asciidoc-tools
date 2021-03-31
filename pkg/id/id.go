package id

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var bracket, invalid, repeat *regexp.Regexp
var prefix string

// Init the regexp
func Init(s string) error {
	bracket = regexp.MustCompile(`(\(.*\))|(\[.*\])|(\{.*\})|(<.*>)`)
	invalid = regexp.MustCompile(`[^a-zA-Z0-9-]`)
	repeat = regexp.MustCompile(`-+`)

	if len(s) == 0 {
		return errors.New("id: empty prefix")
	}

	err := fmt.Errorf("id: invalid prefix %s", s)
	if len(invalid.FindString(s)) != 0 {
		return err
	}
	if s[0] < 'a' || s[0] > 'z' {
		return err
	}
	prefix = s
	return nil
}

// Identify normalizes the ID, it must call after Init()
func Identify(s string) (string, error) {
	if len(prefix) == 0 {
		return prefix, fmt.Errorf("id: Init() not called or failed to initialize")
	}

	s = bracket.ReplaceAllString(s, "")
	s = invalid.ReplaceAllString(s, "-")

	if len(s) == 0 {
		return prefix, nil
	}

	s = strings.ToLower(s)
	if s[0] < 'a' || s[0] > 'z' {
		s = prefix + "-" + s
	}

	s = repeat.ReplaceAllString(s, "-") // handle duplicate '-'

	l := len(s)
	if l > 1 && s[l-1] == '-' { // handle '-' at the end of the ID
		s = s[:l-1]
	}
	return s, nil
}
