package adoc

import (
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var inc, re *regexp.Regexp
var contents = make(map[string]int)

func write(filename string, n int, s string, lines []string) error {
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

func find(filename string) error {
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
			if err = os.Chdir(dir); err != nil {
				return err
			}
			err = find(name)
			if err != nil {
				return err
			}
			if err = os.Chdir(".."); err != nil {
				return err
			}
			continue
		}

		s = re.FindString(line)
		if len(s) != 0 {
			identify := s[2 : len(s)-2]
			if _, ok := contents[identify]; !ok {
				contents[identify] = 1
				continue
			}

			contents[identify]++
			t := "[[" + identify + "-" + strconv.Itoa(contents[identify]) + "]]"
			if err = write(filename, n, t, lines); err != nil {
				return err
			}
			// find this file again for the contents has changed
			return find(filename)
		}
	}
	return nil
}

// Find duplicate IDs in the book and resolve conflicts
func Find(filename string) error {
	inc = regexp.MustCompile(`^include::.*\.adoc\[\]$`)
	re = regexp.MustCompile(`^\[\[.*\]\]$`)

	err := find(filename)
	return err
}
