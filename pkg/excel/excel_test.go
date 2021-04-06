package excel

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// TODO
	Convert(os.Stdout, "testdata/test.xlsx", "Sheet1")
}
