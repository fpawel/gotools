package logfile

import (
	"bytes"
	"fmt"
	"github.com/powerman/structlog"
	"io"
	"os"
	"time"
)

func NewOutput(filenameSuffix string) (io.WriteCloser, error) {
	if err := ensureDir(); err != nil {
		return nil, err
	}
	filename := filename(daytime(time.Now()), filenameSuffix)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0666)
	return &output{f: f}, err
}

type output struct {
	f *os.File
	b bytes.Buffer
}

func (x *output) Close() error {
	return x.f.Close()
}

func (x *output) Write(p []byte) (int, error) {
	x.b.Write(p)
	if bytes.HasSuffix(p, []byte("\n")) {
		if err := x.write(p); err != nil {
			log.PrintErr(err, "line", fmt.Sprintf("%q", x.b.String()), structlog.KeyStack, structlog.Auto)
		}
		x.b.Reset()
	}
	return len(p), nil
}

func (x *output) write(p []byte) error {
	if _, err := fmt.Fprint(x.f, time.Now().Format(layoutDatetime), " "); err != nil {
		return err
	}
	if _, err := x.b.WriteTo(x.f); err != nil {
		return err
	}
	if _, err := x.f.Write(p); err != nil {
		return err
	}
	return nil
}
