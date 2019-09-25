package rungo

import (
	"bytes"
	"fmt"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/gohelp/copydata"
	"github.com/fpawel/gotools/pkg/ccolor"
	"github.com/powerman/structlog"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

func (c Cmd) Exec() {

	log.Info(fmt.Sprintf("run command: %+v", c))

	logFileOutput := NewLogFileOutput()
	defer log.ErrIfFail(logFileOutput.Close)

	writers := []io.Writer{logFileOutput}

	var w *copydata.NotifyWindow
	if c.UseGUI {
		writer := NewNotifyGUIWriter(c.NotifyGUI.WindowClass, c.NotifyGUI.MsgCodeConsole)
		writers = append(writers, writer)
		w = writer.(notifyWriter).w
	}

	cmd := exec.Command(c.ExeName, strings.Fields(c.ExeArgs)...)
	cmd.Stderr = io.MultiWriter(append(writers, ccolor.NewWriter(os.Stderr))...)
	cmd.Stdout = io.MultiWriter(append(writers, ccolor.NewWriter(os.Stdout))...)

	err := c.exec(cmd)
	if err == nil {
		return
	}
	log.PrintErr(err)
	if w != nil {
		go w.NotifyStr(c.NotifyGUI.MsgCodePanic, err.Error())
	}
}

func NewNotifyGUIWriter(windowClass string, msgCodeConsole uintptr) io.Writer {
	s := fmt.Sprintf("%s%d", os.Args[0], time.Now().Unix())
	return notifyWriter{
		w: copydata.NewNotifyWindow(s, windowClass),
		c: msgCodeConsole,
	}
}

func NewLogFileOutput() io.WriteCloser {
	f, err := os.OpenFile(logFileName(), os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &logFileOutput{f: f}
}

func (c Cmd) exec(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

type logFileOutput struct {
	f *os.File
	b bytes.Buffer
}

func (x *logFileOutput) Close() error {
	return x.f.Close()
}

func (x *logFileOutput) Write(p []byte) (int, error) {
	if !bytes.HasSuffix(p, []byte("\n")) {
		x.b.Write(p)
	} else {
		if _, err := fmt.Fprint(x.f, time.Now().Format("15:04:05.000"), " "); err != nil {
			log.PrintErr(err)
		}
		if _, err := x.b.WriteTo(x.f); err != nil {
			log.PrintErr(err)
		}
		if _, err := x.f.Write(p); err != nil {
			log.PrintErr(err)
		}
	}
	return len(p), nil
}

func logFileName() string {
	exeDir := filepath.Dir(os.Args[0])
	t := time.Now()
	logDir := filepath.Join(exeDir, "logs")

	if err := gohelp.EnsuredDir(logDir); err != nil {
		panic(err)
	}

	return filepath.Join(logDir, fmt.Sprintf("%s.log", t.Format("2006-01-02")))
}

type notifyWriter struct {
	w *copydata.NotifyWindow
	c uintptr
}

func (x notifyWriter) Write(p []byte) (int, error) {
	go x.w.NotifyStr(x.c, string(p))
	return len(p), nil
}

var log = structlog.New()
