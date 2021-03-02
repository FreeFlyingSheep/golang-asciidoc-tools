package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/FreeFlyingSheep/golang-asciidoc-tools/pkg/adoc"
)

func main() {
	filename := flag.String("f", "", "Path to the book")
	mode := flag.String("m", "id", "Mode; "+
		"\"id\" - find duplicate IDs in the book and resolve conflicts; "+
		"\"figre\" - generate a list of figures; "+
		"\"table\" - generate a list of tables")
	flag.Parse()

	if len(*filename) == 0 {
		log.Fatalf("You must specifiy a filename by \"-f\" option.")
	}

	name := *filename
	pos := strings.LastIndex(name, "/")
	if pos != -1 {
		dir := name[:pos]
		name = name[pos+1:]
		if err := os.Chdir(dir); err != nil {
			log.Fatalln(err)
		}
	}

	err := adoc.Find(name, *mode)
	if err != nil {
		log.Fatalln(err)
	}
}
