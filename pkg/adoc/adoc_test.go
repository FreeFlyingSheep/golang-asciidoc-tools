package adoc

import (
	"strings"
	"testing"
)

func TestAdoc(t *testing.T) {
	var tests = []struct {
		input string
		mode  string
		ouput string
		wants map[string][]string
	}{
		{
			"test.adoc",
			ID,
			"test",
			map[string][]string{
				"test/duplicate/duplicate.adoc": {
					"duplicate-2",
					"duplicate-3",
				},
			},
		},
		{
			"test.adoc",
			FIGURE,
			"list-of-figures.adoc",
			map[string][]string{
				"list-of-figures.adoc": {
					"List of Figures",
					"Figure",
					"figure",
				},
			},
		},
		{
			"test.adoc",
			TABLE,
			"test/list-of-tables.adoc",
			map[string][]string{
				"test/list-of-tables.adoc": {
					"List of Tables",
					"table-with-cols",
					"Table with `cols`",
					"table-without-cols",
					"Table without `cols`",
				},
			},
		},
	}

	for _, test := range tests {
		path := "testdata/" + test.input
		if err := Find(path, test.mode, test.ouput); err != nil {
			t.Error(err)
		}

		got, err := List()
		if err != nil {
			t.Error(err)
		}

		for file, want := range test.wants {
			var name string
			switch test.mode {
			case ID:
				name = "testdata/" + file
			case FIGURE:
				fallthrough
			case TABLE:
				name = file
			}

			if _, ok := got[name]; !ok {
				t.Errorf("List() not contains file: %s", file)
				continue
			}
			var lines string
			for _, line := range got[name] {
				lines += line + "\n"
			}
			for _, line := range want {
				if !strings.Contains(lines, line) {
					t.Errorf("List() not contains %s", line)
				}
			}
		}
	}
}
