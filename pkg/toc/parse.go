package toc

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

const (
	// MAXLEVEL means the maximum number of levels of the book
	MAXLEVEL = 5
	// MAXSECTION The maximum number of sections in the book
	MAXSECTION = 100000
)

// Section partition the book into a content hierarchy
type Section struct {
	Number []int
	Title  string
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

func parseNum(s string, sep string) ([]int, error) {
	var nums []int
	parts := strings.Split(s, sep)
	if len(parts) > MAXLEVEL { // level > MaxLevel
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
func Parse(body []byte, NumSep, TitleSep string, book string) (*TOC, error) {
	bookSection := &Section{
		Number: []int{},
		Title:  book,
	}
	toc := &TOC{
		Sections: []*Section{bookSection},
	}

	lines := strings.Split(string(body), "\n")
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

		num, err := parseNum(line[:pos], NumSep)
		if err != nil {
			return nil, err
		} else if num == nil { // just ignore it
			continue
		}
		level := len(num)
		if level > toc.Level {
			toc.Level = level
		}

		title := strings.Trim(line[pos:], " ")
		section := &Section{
			Number: num,
			Title:  title,
		}
		toc.Sections = append(toc.Sections, section)
	}

	sort.Sort(toc)
	return toc, nil
}
