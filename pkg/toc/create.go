package toc

import (
	"errors"
	"fmt"
	"os"

	"github.com/FreeFlyingSheep/golang-asciidoc-tools/pkg/id"
)

// MaxName means the maximum length of the filename
const MaxName = 255

func write(n int, identify, s string) error {
	var err error

	name := identify + ".adoc"
	if len(name) > MaxName {
		return errors.New("toc: filename too long")
	}
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(s)
	return err
}

func create(n int, sections []*Section) (int, error) {
	var err error
	identify := id.Identify(sections[n].Title)
	level := len(sections[n].Number) + 1
	var symbol string
	for i := 0; i < level; i++ {
		symbol += "="
	}

	var contents string
	if n != 0 {
		contents = fmt.Sprintf("[[%s]]\n", identify)
	}
	contents = fmt.Sprintf("%s%s %s\n", contents, symbol, sections[n].Title)

	i := n + 1
	first := true
	for i < len(sections) {
		length := len(sections[i].Number)
		if length == level {
			path := identify + "/" + id.Identify(sections[i].Title)
			contents = fmt.Sprintf("%s\ninclude::%s.adoc[]\n", contents, path)
			if first {
				if err = os.Mkdir(identify, os.ModeDir); err != nil {
					return i, err
				}
				if err = os.Chdir(identify); err != nil {
					return i, err
				}
				first = false
			}
		} else if length < level {
			break
		}

		i, err = create(i, sections)
		if err != nil {
			return i, err
		}
	}

	if !first {
		if err = os.Chdir(".."); err != nil {
			return i, err
		}
	}
	err = write(n, identify, contents)
	return i, err
}

// Create files via the table of contents
func Create(toc *TOC, output string) error {
	_, err := os.Stat(output)
	if err == nil {
		return fmt.Errorf("toc: %s already exists", output)
	}

	if err = os.MkdirAll(output, os.ModeDir); err != nil {
		return err
	}
	if err = os.Chdir(output); err != nil {
		return err
	}

	id.Init()
	n, err := create(0, toc.Sections)
	if err != nil {
		return err
	}

	if n != len(toc.Sections) {
		return errors.New("toc: create files error")
	}
	return nil
}
