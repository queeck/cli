package selected

import (
	"fmt"
	"strconv"

	"github.com/queeck/cli/internal/pkg/styles"
)

const (
	padShift = 2
)

type printerDefault struct {
	//
}

func PrinterDefault() Printer {
	return &printerDefault{}
}

func (pd *printerDefault) Selected(code string, pad int) string {
	return styles.ColorForegroundHighlight(fmt.Sprintf("> %-"+strconv.Itoa(pad+padShift)+"s", code))
}

func (pd *printerDefault) Unselected(code string, pad int) string {
	return styles.ColorForegroundSubtle(fmt.Sprintf("  %-"+strconv.Itoa(pad+padShift)+"s", code))
}
