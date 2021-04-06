package main

import (
	"flag"
	"log"
	"os"

	"github.com/FreeFlyingSheep/golang-asciidoc-tools/pkg/excel"
)

func main() {
	filename := flag.String("f", "", "Path to the file")
	sheet := flag.String("s", "", "Sheet name")
	ouput := flag.String("o", "", "output file")
	flag.Parse()

	if len(*filename) == 0 {
		log.Fatalf("You must specifiy a filename by \"-f\" option.")
	}

	if len(*sheet) == 0 {
		log.Fatalf("You must specifiy a sheet by \"-s\" option.")
	}

	if len(*ouput) == 0 {
		if err := excel.Convert(os.Stdout, *filename, *sheet); err != nil {
			log.Fatalln(err)
		}
	} else {
		file, err := os.Create(*ouput)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		if err := excel.Convert(file, *filename, *sheet); err != nil {
			log.Fatalln(err)
		}
	}
}
