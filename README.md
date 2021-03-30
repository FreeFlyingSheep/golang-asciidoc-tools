# Golang AsciiDoc Tools

Golang-based tools for AsciiDoc.

## Book

### Features of Book

- Parse a table of contents (TOC) in a text file.
- Create files via the TOC.

### Examples of Book

```bash
go build cmd/book/main.go
./main.exe -c -f toc.txt -b name -o book
```

## List

### Features of List

- Find duplicate IDs in the book and resolve conflicts.
- Generate a list of figures.
- Generate a list of tables.

### Examples of List

```bash
go build cmd/list/main.go
./main.exe -f book.adoc
```
