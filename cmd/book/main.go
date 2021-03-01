package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/FreeFlyingSheep/golang-asciidoc-tools/pkg/toc"
)

func main() {
	numSep := flag.String("n", ".",
		"Number separator, e.g., \".\" is a separator of \"1.1\"")
	titleSep := flag.String("t", " ",
		"Title separator, e.g., \" \" is a separator of \"1.1 xxx\"")
	filename := flag.String("f", "",
		"A text file that contains the table of contents")
	ouput := flag.String("o", "",
		"Output directory, if not specified, just print the table of contents")
	book := flag.String("b", "", "Book name")
	flag.Parse()

	if len(*filename) == 0 {
		log.Fatalf("You must specifiy a filename by \"-f\" option.")
	}

	if len(*book) == 0 {
		log.Fatalf("You must specifiy a book name by \"-b\" option.")
	}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	t, err := toc.Parse(body, *numSep, *titleSep, *book)
	if err != nil {
		log.Fatalln(err)
	}

	if len(*ouput) == 0 {
		toc.Print(t)
	} else {
		err = toc.Create(t, *ouput)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
