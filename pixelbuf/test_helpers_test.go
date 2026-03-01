package pixelbuf

func solidBuffer(w, h int, c Color) *Buffer {
	buf := NewBuffer(w, h)
	buf.Clear(c)
	return buf
}

