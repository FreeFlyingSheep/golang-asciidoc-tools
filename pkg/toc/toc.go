package toc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/FreeFlyingSheep/golang-asciidoc-tools/pkg/id"
)

const (
	// MAXLEVEL means the maximum number of levels of the book
	MAXLEVEL = 5
	// MAXSECTION means the maximum number of sections in the book
	MAXSECTION = 100000
	// MAXNAME means the maximum length of the filename
	MAXNAME = 255
)

// Section partition the book into a content hierarchy
type Section struct {
	Number   []int
	Level    int
	Title    string
	Identify string
	Path     string
	Content  string
}

// TOC means table of contents
type TOC struct {
	Level    int
	Total    int
	Sections []*Section
}

func (x TOC) Len() int {
	return len(x.Sections)
}

func (x TOC) Less(i, j int) bool {
	num1 := x.Sections[i].Number
	num2 := x.Sections[j].Number
	len1 := len(num1)
	len2 := len(num2)

	k := 0
	for ; k < len1 && k < len2; k++ {
		if num1[k] != num2[k] {
			return num1[k] < num2[k] // e.g. [1 1] < [1 2]
		}
	}
	return len1 < len2 // e.g. [1] < [1 1] (k = 1)
}

func (x TOC) Swap(i, j int) {
	x.Sections[i], x.Sections[j] = x.Sections[j], x.Sections[i]
}

func parseNum(s string, sep string, l int) ([]int, error) {
	var nums []int
	parts := strings.Split(s, sep)
	if len(parts) > l { // level > l
		return nil, nil
	}

	for _, part := range parts {
		if len(part) == 0 { // empty line
			continue
		}

		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	return nums, nil
}

// Parse all sections from the text
func Parse(body, NumSep, TitleSep, book string, level int) (*TOC, error) {
	bookSection := &Section{
		Number: []int{},
		Title:  book,
	}
	toc := &TOC{
		Sections: []*Section{bookSection},
	}

	lines := strings.Split(body, "\n")
	toc.Total = len(lines)
	if toc.Total > MAXSECTION {
		return nil, errors.New("toc: too many sections")
	}

	for _, line := range lines {
		line = strings.Trim(line, " ")
		if len(line) == 0 { // empty line
			toc.Total--
			continue
		}

		pos := strings.Index(line, TitleSep)
		if pos == -1 {
			return nil, errors.New("toc: section number not found")
		}

		if level == 0 || level > MAXLEVEL {
			level = MAXLEVEL
		}
		num, err := parseNum(line[:pos], NumSep, level)
		if err != nil {
			return nil, err
		} else if num == nil { // just ignore it
			continue
		}
		l := len(num)
		if l > toc.Level {
			toc.Level = l
		}

		title := strings.Trim(line[pos:], " ")
		section := &Section{
			Number: num,
			Level:  l,
			Title:  title,
		}
		toc.Sections = append(toc.Sections, section)
	}

	sort.Sort(toc)
	return toc, nil
}

// Write the TOC
func Write(out io.Writer, toc *TOC) {
	fmt.Fprintf(out, "Level: %d\n", toc.Level)
	fmt.Fprintf(out, "Total: %d\n", toc.Total)

	for i, section := range toc.Sections {
		length := len(section.Number) - 1
		for j, num := range section.Number {
			fmt.Fprintf(out, "%d", num)
			if j != length {
				fmt.Fprintf(out, ".")
			}
		}

		if i != 0 {
			fmt.Fprintf(out, " ")
		}
		fmt.Fprintf(out, "%s\n", section.Title)
	}
}

func makeHeader(section *Section, toc *TOC, custom bool) error {
	identify, err := id.Identify(section.Title)
	if err != nil {
		return err
	}
	section.Identify = identify

	var content string
	if custom {
		content = fmt.Sprintf("[[%s]]\n", section.Identify)
	}

	var symbol string
	for i := 0; i < section.Level+1; i++ {
		symbol += "="
	}
	content = fmt.Sprintf("%s%s %s\n", content, symbol, section.Title)
	section.Content = content
	return nil
}

func generate(toc *TOC, custom bool) error {
	section := toc.Sections[0]
	if err := makeHeader(section, toc, false); err != nil {
		return err
	}
	parent := section

	for i := 1; i < len(toc.Sections); i++ {
		section = toc.Sections[i]
		last := toc.Sections[i-1]
		if err := makeHeader(section, toc, custom); err != nil {
			return err
		}

		if section.Level > last.Level+1 {
			return errors.New("toc: missing parent section")
		} else if section.Level == last.Level+1 { // enter the subdirectory
			parent = last
		} else if section.Level < last.Level { // back to the parent directory
			n := i
			for n > 0 && section.Level <= toc.Sections[n].Level {
				n--
			}
			parent = toc.Sections[n]
		}

		section.Path = parent.Path + "/" + parent.Identify
		name := parent.Identify + "/" + section.Identify
		parent.Content = fmt.Sprintf("%s\ninclude::%s.adoc[]\n", parent.Content, name)
	}
	return nil
}

// Generate files via the table of contents
func Generate(toc *TOC, custom bool, prefix, output string) error {
	_, err := os.Stat(output)
	if err == nil {
		return fmt.Errorf("toc: %s already exists", output)
	}

	if err := id.Init(prefix); err != nil {
		return err
	}
	toc.Sections[0].Path = output // use Path as file directory for now
	err = generate(toc, custom)
	if err != nil {
		return err
	}

	for _, section := range toc.Sections { // add filename to Path
		section.Path += "/" + section.Identify + ".adoc"
	}
	return nil
}
