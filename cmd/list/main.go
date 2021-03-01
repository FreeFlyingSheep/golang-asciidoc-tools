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

	err := adoc.Find(name)
	if err != nil {
		log.Fatalln(err)
	}
}
