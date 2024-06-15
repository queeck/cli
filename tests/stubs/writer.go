package stubs

import (
	"io"

	"github.com/queeck/cli/tests/utils"
)

type WriterWithFrames interface {
	io.Writer
	FramesWithLF() []string
}

type writer struct {
	written []byte
}

func Writer() WriterWithFrames {
	return &writer{written: make([]byte, 0)}
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.written = append(w.written, p...)
	return len(p), nil
}

func (w *writer) FramesWithLF() []string {
	return utils.Frames(utils.ToLF(w.written))
}
