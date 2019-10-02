package copydata

import "io"

func NewWriter(msg uintptr, srcWndClass, destWndClass string) io.Writer {
	return guiWriter{
		msg:          msg,
		srcWndClass:  srcWndClass,
		destWndClass: destWndClass,
	}
}

type guiWriter struct {
	msg                       uintptr
	srcWndClass, destWndClass string
}

func (x guiWriter) Write(p []byte) (int, error) {
	WndClass{x.srcWndClass, x.destWndClass}.SendString(x.msg, string(p))
	return len(p), nil
}
