package adoc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Modes
const (
	ID     = "id"
	FIGURE = "figure"
	TABLE  = "table"
)

var inc, id, title, figure, table *regexp.Regexp
var ids = make(map[string]int)
var contents = make(map[string][]string)
var list []string
var mode string
var out string

func countID(line string) string {
	var t string
	s := id.FindString(line)
	if len(s) != 0 {
		identify := s[2 : len(s)-2]
		if _, ok := ids[identify]; !ok {
			ids[identify] = 1
			return t
		}

		ids[identify]++
		t = "[[" + identify + "-" + strconv.Itoa(ids[identify]) + "]]"
	}
	return t
}

func findID(filename string, n int, s string, lines []string) {
	if _, ok := contents[filename]; !ok {
		contents[filename] = lines
	}
	contents[filename][n] = s
}

func findList(s1, s2, s3 string) {
	r := id.FindString(s1)
	if len(r) == 0 {
		return
	}
	s := title.FindString(s2)
	if len(s) == 0 {
		return
	}

	var t string
	switch mode {
	case FIGURE:
		t = figure.FindString(s3)
	case TABLE:
		t = table.FindString(s3)
	}
	if len(t) == 0 {
		return
	}

	// len("[[") == 2, len("]]") == 2, len(".") == 1
	t = "* <<" + r[2:len(r)-2] + "," + s[1:] + ">>"
	list = append(list, t)
}

func List() (map[string][]string, error) {
	if len(mode) == 0 {
		return nil, errors.New("adoc: Find() not called or failed to find")
	}

	switch mode {
	case ID:
		return contents, nil
	case FIGURE:
		contents[out] = []string{"== List of Figures\n"}
	case TABLE:
		contents[out] = []string{"== List of Tables\n"}
	}
	l := len(list)
	if l > 0 {
		list[l-1] += "\n"
	}
	contents[out] = append(contents[out], list...)
	return contents, nil
}

func find(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(body), "\n")
	for n, line := range lines {
		s := inc.FindString(line)
		if len(s) != 0 {
			p := filepath.Dir(path)
			p = filepath.ToSlash(p)
			// len("include::") == 9, len("[]") == 2
			p += "/" + s[9:len(s)-2]
			err = find(p)
			if err != nil {
				return err
			}
			continue
		}

		switch mode {
		case ID:
			s = countID(line)
			if len(s) != 0 {
				findID(path, n, s, lines)
			}
		case FIGURE:
			fallthrough
		case TABLE:
			if n >= 2 {
				findList(lines[n-2], lines[n-1], line)
			}
		}
	}
	return nil
}

// Find things in the book
func Find(path, findMode, ouput string) error {
	mode = findMode
	out = ouput

	ids = make(map[string]int)
	contents = make(map[string][]string)
	list = []string{}

	inc = regexp.MustCompile(`^include::.*\.adoc\[\]$`)
	id = regexp.MustCompile(`^\[\[.*\]\]$`)
	switch mode {
	case ID:
		// do nothing
	case FIGURE:
		title = regexp.MustCompile(`^\.\S.*$`)
		figure = regexp.MustCompile(`^image::.*\[\]$`)
	case TABLE:
		title = regexp.MustCompile(`^\..*$`)
		table = regexp.MustCompile(`(^\|===.*$)|(^\[.*\]$)`)
	default:
		return fmt.Errorf("find: unknown type %s", mode)
	}
	return find(path)
}
