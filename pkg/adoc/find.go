package adoc

import (
	"fmt"
	"io/ioutil"
	"os"
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
var contents = make(map[string]int)
var list []string

func findID(line string) string {
	var t string
	s := id.FindString(line)
	if len(s) != 0 {
		identify := s[2 : len(s)-2]
		if _, ok := contents[identify]; !ok {
			contents[identify] = 1
			return t
		}

		contents[identify]++
		t = "[[" + identify + "-" + strconv.Itoa(contents[identify]) + "]]"
	}
	return t
}

func findList(s1, s2, s3 string, mode string) {
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

func writeID(filename string, n int, s string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for i, line := range lines {
		if i == n {
			_, err := file.WriteString(s + "\n")
			if err != nil {
				return err
			}
		} else {
			_, err := file.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeList(mode string) error {
	file, err := os.Create("table-of-contents.adoc")
	if err != nil {
		return err
	}
	defer file.Close()
	switch mode {
	case FIGURE:
		_, err = file.WriteString("== List of figures\n\n")
	case TABLE:
		_, err = file.WriteString("== List of tables\n\n")
	}
	if err != nil {
		return err
	}

	for _, item := range list {
		_, err := file.WriteString(item + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func find(filename, mode string) error {
	file, err := os.Open(filename)
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
			pos := strings.LastIndex(s, "/")
			dir := s[9:pos]             // len("include::") == 9
			name := s[pos+1 : len(s)-2] // len("[]") == 2

			curr, err := os.Getwd()
			if err != nil {
				return err
			}
			if err = os.Chdir(dir); err != nil {
				return err
			}
			err = find(name, mode)
			if err != nil {
				return err
			}
			if err = os.Chdir(curr); err != nil {
				return err
			}
			continue
		}

		switch mode {
		case ID:
			s = findID(line)
			if len(s) != 0 {
				if err = writeID(filename, n, s, lines); err != nil {
					return err
				}
				// find this file again for the contents has changed
				return find(filename, mode)
			}
		case FIGURE:
			fallthrough
		case TABLE:
			if n >= 2 {
				findList(lines[n-2], lines[n-1], line, mode)
			}

		}
	}
	return nil
}

// Find duplicate IDs in the book and resolve conflicts
func Find(filename, mode string) error {
	inc = regexp.MustCompile(`^include::.*\.adoc\[\]$`)
	id = regexp.MustCompile(`^\[\[.*\]\]$`)
	switch mode {
	case ID:
		return find(filename, mode)
	case FIGURE:
		title = regexp.MustCompile(`^\..*$`)
		figure = regexp.MustCompile(`^image::.*\[\]$`)
	case TABLE:
		title = regexp.MustCompile(`^\..*$`)
		table = regexp.MustCompile(`(^\|===.*$)|(^\[.*\]$)`)
	default:
		return fmt.Errorf("find: unknown type %s", mode)
	}

	if err := find(filename, mode); err != nil {
		return err
	}
	return writeList(mode)
}
