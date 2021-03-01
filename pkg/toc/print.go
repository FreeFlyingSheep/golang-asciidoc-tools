package toc

import "fmt"

// Print the TOC
func Print(toc *TOC) {
	fmt.Printf("Level: %d\n", toc.Level)
	fmt.Printf("Total: %d\n", toc.Total)

	for i, section := range toc.Sections {
		length := len(section.Number) - 1
		for j, num := range section.Number {
			fmt.Printf("%d", num)
			if j != length {
				fmt.Printf(".")
			}
		}

		if i != 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%s\n", section.Title)
	}
}
