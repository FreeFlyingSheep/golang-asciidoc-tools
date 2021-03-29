package id

import (
	"testing"
)

func TestInit(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"", "id: empty prefix"},
		{"123", "id: invalid prefix 123"},
		{"abc", ""},
		{"DEF", "id: invalid prefix DEF"},
		{"a(0)", "id: invalid prefix a(0)"},
	}

	for _, test := range tests {
		got := Init(test.input)
		var res string
		if got != nil {
			res = got.Error()
		}
		if res != test.want {
			t.Errorf("Init(%s) = %v", test.input, got)
		}
	}
}

func TestIdentify(t *testing.T) {
	prefix := "section"
	var tests = []struct {
		input string
		want  string
	}{
		{"", prefix},
		{"abc", "abc"},
		{"Def", "def"},
		{"123", prefix + "-123"},
		{"{a}", prefix},
		{"[b1]", prefix},
		{"(c)2", prefix + "-2"},
		{"<d>EF", "ef"},
		{"+-*/", prefix},
		{"abc-(123)-DEF-+-", "abc-def"},
	}

	if err := Init(prefix); err != nil {
		t.Errorf("Init(%s) = %v", prefix, err)
	}
	for _, test := range tests {
		got, err := Identify(test.input)
		if err != nil {
			t.Errorf("Identify(%s) = %v", prefix, got)
		}
		if got != test.want {
			t.Errorf("Identify(%s) = %v", test.input, got)
		}
	}
}
