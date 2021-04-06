package excel

import (
	"bytes"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	var tests = []struct {
		input string
		sheet string
		wants []string
	}{
		{"test.xlsx", "Sheet1",
			[]string{
				"[cols=\"5*1\"]",
				"|===",
				"|A1",
				"|B1",
				"|C1",
				"|D1",
				"|E1",
				"2+|A2",
				"|C2",
				".3+|D2",
				"|E2",
				"|A3",
				"|B3",
				"|C3",
				"|E3",
				"|A4",
				"|B4",
				"|C4",
				"|E4",
				"|A5",
				"2.2+|B5",
				"|D5",
				"|E5",
				"|A6",
				"|D6",
				"|E6",
				"|A7",
				"|B7",
				"|C7",
				"|D7",
				"|E7",
			},
		},
	}

	out := new(bytes.Buffer)
	for _, test := range tests {
		if err := Convert(out, "testdata/"+test.input, test.sheet); err != nil {
			t.Error(err)
		}

		got := out.String()
		for _, want := range test.wants {
			if !strings.Contains(got, want) {
				t.Errorf("Convert() not contains %s", want)
			}
		}
	}
}
