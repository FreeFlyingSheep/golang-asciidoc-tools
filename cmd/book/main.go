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
	prefix := flag.String("p", "section",
		"The prefix added when ID is invalid")
	ouput := flag.String("o", "",
		"Output directory, if not specified, just print the table of contents")
	book := flag.String("b", "", "Book name")
	level := flag.Int("l", 0,
		"The maximum number of levels of the book, \"0\" means no limit")
	custom := flag.Bool("i", false, "Generate custom section IDs")
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

	t, err := toc.Parse(body, *numSep, *titleSep, *book, *level)
	if err != nil {
		log.Fatalln(err)
	}

	if len(*ouput) == 0 {
		toc.Write(os.Stdout, t)
	} else {
		err = toc.Create(t, *custom, *prefix, *ouput)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
