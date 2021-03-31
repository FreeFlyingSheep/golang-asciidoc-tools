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
		custom   bool
		prefix   string
		path     string
		max      int
		total    int
		sections []string
		files    []string
	}{
		{
			"1{titleSep}title\n" +
				"1{numSep}1{titleSep}1.1\n" +
				"1{numSep}1{numSep}1{titleSep}1.1.1\n" +
				"1{numSep}2{numSep}1{titleSep}world\n" +
				"2{titleSep}go\n" +
				"1{numSep}1{numSep}1{numSep}1{titleSep}1.1.1.1\n" +
				"1{numSep}2{titleSep}hello\n",
			".",
			" ",
			"book",
			2,
			true,
			"section",
			"test",
			2,
			7,
			[]string{"book", "title", "1.1", "hello", "go"},
			[]string{
				"test/book.adoc:",
				"= book",
				"include::book/title.adoc[]",
				"include::book/go.adoc[]",
				"test/book/title.adoc:",
				"[[title]]",
				"== title",
				"include::title/section-1-1.adoc[]",
				"include::title/hello.adoc[]",
				"test/book/title/section-1-1.adoc:",
				"[[section-1-1]]",
				"=== 1.1",
				"test/book/title/hello.adoc:",
				"[[hello]]",
				"=== hello",
				"test/book/go.adoc:",
				"[[go]]",
				"== go",
			},
		},
	}

	numSep := regexp.MustCompile(`\{numSep\}`)
	titleSep := regexp.MustCompile(`\{titleSep\}`)
	for _, test := range tests {
		input := numSep.ReplaceAllString(test.input, test.numSep)
		input = titleSep.ReplaceAllString(input, test.titleSep)

		toc, err := Parse(input, test.numSep, test.titleSep, test.title, test.level)
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
		for _, s := range test.sections {
			if !strings.Contains(got, s) {
				t.Errorf("Write(toc) not contains %s", s)
			}
		}

		err = Generate(toc, test.custom, test.prefix, test.path)
		if err != nil {
			t.Error(err)
		}
		var files string
		for _, section := range toc.Sections {
			files += section.Path + ":\n"
			files += section.Content + "\n"
		}
		for _, s := range test.files {
			if !strings.Contains(files, s) {
				t.Errorf("Generate(toc) not contains %s", s)
			}
		}
	}
}
