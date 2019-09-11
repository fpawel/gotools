package rungo

import (
	"fmt"
	"github.com/fpawel/gohelp/copydata"
	"os"
	"time"
)

func (c NotifyGUI) newWriter() notifyWriter {
	s := fmt.Sprintf("%s%d", os.Args[0], time.Now().Unix())
	return notifyWriter{
		w: copydata.NewNotifyWindow(s, c.WindowClass),
		c: c.MsgCodeConsole,
	}
}

type notifyWriter struct {
	w *copydata.NotifyWindow
	c uintptr
}

func (x notifyWriter) Write(p []byte) (int, error) {
	go x.w.NotifyStr(x.c, string(p))
	return len(p), nil
}
