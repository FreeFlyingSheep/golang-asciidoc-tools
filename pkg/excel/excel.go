package excel

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type tableCell struct {
	exist bool
	col   int
	row   int
	value string
}

var table [][]tableCell
var merge []string
var cols int

func read(input, sheet string) error {
	table = [][]tableCell{}

	f, err := excelize.OpenFile(input)
	if err != nil {
		return err
	}

	cells, err := f.GetMergeCells(sheet)
	if err != nil {
		return err
	}
	for _, cell := range cells {
		merge = append(merge, cell[0])
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return errors.New("excel: no rows")
	}
	cols = len(rows[0])

	for _, row := range rows {
		cells := []tableCell{}
		for _, cell := range row {
			cells = append(cells, tableCell{true, 1, 1, cell})
		}
		table = append(table, cells)
	}
	return nil
}

func mergeCell(cell string) error {
	pos := strings.Index(cell, ":")
	if pos == -1 {
		return fmt.Errorf("excel: MergeCells %s not contain \":\"", cell)
	}

	hcell := cell[:pos]
	hcol, hrow, err := excelize.CellNameToCoordinates(hcell)
	if err != nil {
		return err
	}

	vcell := cell[pos+1:]
	vcol, vrow, err := excelize.CellNameToCoordinates(vcell)
	if err != nil {
		return err
	}

	crow := hrow - 1
	ccol := hcol - 1
	for row := crow; row < vrow; row++ {
		for col := ccol; col < vcol; col++ {
			table[row][col].exist = false
			table[row][col].row = 0
			table[row][col].col = 0
		}
	}
	table[crow][ccol].exist = true
	table[crow][ccol].col = vcol - hcol + 1
	table[crow][ccol].row = vrow - hrow + 1
	return nil
}

func write(out io.Writer) {
	fmt.Fprintf(out, "[cols=\"%d*1\"]\n|===", cols)
	for _, row := range table {
		fmt.Fprintf(out, "\n")
		for _, cell := range row {
			if !cell.exist {
				continue
			}

			isMerge := false
			if cell.col > 1 {
				fmt.Fprintf(out, "%d", cell.col)
				isMerge = true
			}
			if cell.row > 1 {
				fmt.Fprintf(out, ".%d", cell.row)
				isMerge = true
			}
			if isMerge {
				fmt.Fprintf(out, "+")
			}
			fmt.Fprintf(out, "|%s\n", cell.value)
		}
	}
	fmt.Fprintf(out, "|===\n")
}

// Convert a Excel table to a AsciiDoc table
func Convert(out io.Writer, input, sheet string) error {
	if err := read(input, sheet); err != nil {
		return err
	}

	for _, cell := range merge {
		if err := mergeCell(cell); err != nil {
			return err
		}
	}

	write(out)
	return nil
}
