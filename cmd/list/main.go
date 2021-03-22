package main

import (
	"flag"
	"log"

	"github.com/FreeFlyingSheep/golang-asciidoc-tools/pkg/adoc"
)

func main() {
	filename := flag.String("f", "", "Path to the book")
	mode := flag.String("m", "id", "Mode; "+
		"\"id\" - find duplicate IDs in the book and resolve conflicts; "+
		"\"figure\" - generate a list of figures; "+
		"\"table\" - generate a list of tables")
	ouput := flag.String("o", "table-of-contents.adoc",
		"Output file, only works when mode is not \"id\"")
	flag.Parse()

	if len(*filename) == 0 {
		log.Fatalf("You must specifiy a filename by \"-f\" option.")
	}

	if len(*ouput) == 0 {
		log.Fatalf("You must specifiy a filename by \"-o\" option.")
	}

	err := adoc.Find(*filename, *mode, *ouput)
	if err != nil {
		log.Fatalln(err)
	}
}
