package logfile

import (
	"bytes"
	"fmt"
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
}

func (x *output) Close() error {
	return x.f.Close()
}

func (x *output) Write(p []byte) (int, error) {
	for _, p := range bytes.Split(p, []byte{'\n'}) {
		if len(p) == 0 {
			continue
		}
		if _, err := fmt.Fprint(x.f, time.Now().Format(layoutDatetime), " "); err != nil {
			return 0, err
		}
		if _, err := x.f.Write(p); err != nil {
			return 0, err
		}
		if !bytes.HasSuffix(p, []byte("\n")) {
			if _, err := x.f.WriteString("\n"); err != nil {
				return 0, err
			}
		}
	}
	return len(p), nil
}
