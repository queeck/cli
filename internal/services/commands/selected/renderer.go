package selected

import (
	"slices"
	"strings"

	"github.com/queeck/cli/internal/services/commands"
)

type Entry struct {
	Code       string
	IsSelected bool
}

type Selection interface {
	Render(command commands.Command, code string) string
}

type Printer interface {
	Selected(code string, pad int) string
	Unselected(code string, pad int) string
}

type Commands struct {
	bus     commands.Bus
	printer Printer
}

func New(bus commands.Bus) *Commands {
	return &Commands{
		bus:     bus,
		printer: PrinterDefault(),
	}
}

func (c *Commands) Render(command commands.Command, code string) string {
	list := command.Commands()
	entriesMatrix := make([][]Entry, 0, len(list))

	for {
		entries := make([]Entry, 0, 1)

		for _, model := range command.Commands() {
			entries = append(entries, Entry{Code: model.Code, IsSelected: model.Code == code})
		}

		entriesMatrix = append(entriesMatrix, entries)

		code = command.Code()

		command = c.bus.Parent(command)
		if command == nil {
			break
		}
	}
	return view(transpose(revert(entriesMatrix)), c.printer)
}

func revert(entriesMatrix [][]Entry) [][]Entry {
	slices.Reverse(entriesMatrix)
	return entriesMatrix
}

func transpose(entriesMatrix [][]Entry) [][]Entry {
	maxSide := len(entriesMatrix)

	for i := range entriesMatrix {
		if maxSide < len(entriesMatrix[i]) {
			maxSide = len(entriesMatrix[i])
		}
	}

	transposed := make([][]Entry, maxSide)
	for i := range maxSide {
		transposed[i] = make([]Entry, maxSide)
	}

	for i := range entriesMatrix {
		for j := range len(entriesMatrix[i]) {
			transposed[j][i] = entriesMatrix[i][j]
		}
	}

	return transposed
}

func view(entries [][]Entry, printer Printer) string {
	lines := make([]string, 0, len(entries))
	for i := range entries {
		keys := make([]string, 0, len(entries[i]))
		for j := range entries[i] {
			entry := entries[i][j]
			if entry.Code != "" {
				printing := printer.Unselected
				if entry.IsSelected {
					printing = printer.Selected
				}
				key := printing(entry.Code, maxKeyLength(entries, j))
				keys = append(keys, key)
			}
		}
		lines = append(lines, strings.Join(keys, " "))
	}
	return strings.Join(lines, "\n")
}

func maxKeyLength(entries [][]Entry, vertical int) int {
	length := 0
	for i := range entries {
		if vertical < len(entries[i]) {
			key := entries[i][vertical].Code
			if len(key) > length {
				length = len(key)
			}
		}
	}
	return length
}
