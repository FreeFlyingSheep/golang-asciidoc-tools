package toc

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestToc(t *testing.T) {
	var tests = []struct {
		input    string
		numSep   string
		titleSep string
		title    string
		level    int
		max      int
		total    int
	}{
		{"1{titleSep}title\n" +
			"1{numSep}1{titleSep}1.1\n" +
			"1{numSep}1{numSep}1{titleSep}1.1.1\n" +
			"1{numSep}1{numSep}1{numSep}1{titleSep}1.1.1.1\n" +
			"1{numSep}1{numSep}1{numSep}1{numSep}1{titleSep}1.1.1.1.1\n" +
			"1{numSep}1{numSep}1{numSep}1{numSep}1{numSep}1{titleSep}1.1.1.1.1.1\n" +
			"1{numSep}2{titleSep}hello\n" +
			"1{numSep}2{numSep}1{titleSep}world\n" +
			"2{titleSep}go\n", ".", " ", "book", 3, 3, 9},
	}

	numSep := regexp.MustCompile(`\{numSep\}`)
	titleSep := regexp.MustCompile(`\{titleSep\}`)
	for _, test := range tests {
		input := numSep.ReplaceAllString(string(test.input), test.numSep)
		input = titleSep.ReplaceAllString(string(input), test.titleSep)

		toc, err := Parse([]byte(input), test.numSep,
			test.titleSep, test.title, test.level)
		if err != nil {
			t.Error(err)
		}
		if toc.Level != test.max || toc.Total != test.total {
			t.Errorf("Parse(%s):\ntoc.Level: %d\ntoc.Total: %d",
				input, toc.Level, toc.Total)
		}

		out := new(bytes.Buffer)
		Write(out, toc)
		got := out.String()
		if !strings.Contains(got, test.title) {
			t.Errorf("Write(toc) not contains %s", test.title)
		}
	}
}
