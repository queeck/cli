package bus

type window struct {
	height int
	width  int
}

func (w *window) Height() int {
	return w.height
}

func (w *window) Width() int {
	return w.width
}
