package rungo

import (
	"fmt"
	"github.com/fpawel/gohelp/copydata"
	"github.com/fpawel/gotools/internal/loggo/data"
	"github.com/fpawel/gotools/pkg/ccolor"
	"github.com/powerman/structlog"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Cmd struct {
	ExeName   string
	ExeArgs   string
	UseGUI    bool
	NotifyGUI NotifyGUI
}

type NotifyGUI struct {
	MsgCodeConsole uintptr
	MsgCodePanic   uintptr
	WindowClass    string
}

func (c Cmd) Exec() error {
	db := data.New(data.KindText, c.ExeName)
	defer log.ErrIfFail(db.Close)
	writers := []io.Writer{db}

	var notifier notifyWriter

	if c.UseGUI {
		notifier = c.NotifyGUI.newWriter()
		writers = append(writers, notifier)
	}

	cmd := exec.Command(c.ExeName, strings.Fields(c.ExeArgs)...)
	cmd.Stderr = io.MultiWriter(append(writers, ccolor.NewWriter(os.Stderr))...)
	cmd.Stdout = io.MultiWriter(append(writers, ccolor.NewWriter(os.Stdout))...)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if c.UseGUI {
			go notifier.w.NotifyStr(c.NotifyGUI.MsgCodePanic, err.Error())
		}
		return log.Err(err)
	}
	return nil
}

var log = structlog.New()

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
