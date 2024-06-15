package stubs

import "io"

func ReaderEmpty() io.Reader {
	return &inputMock{
		ReadFunc: func([]byte) (int, error) {
			return 0, io.EOF
		},
	}
}

func ReaderForText(text string) io.Reader {
	in := []byte(text)
	done := false
	return &inputMock{
		ReadFunc: func(p []byte) (int, error) {
			if done {
				return 0, io.EOF
			}
			for i := range p {
				if i == len(in) {
					done = true
					return len(in), io.EOF
				}
				p[i] = in[i]
			}
			if len(in) > len(p) {
				in = in[len(p):]
				return len(p), nil
			}
			done = true
			return len(in), nil
		},
	}
}
