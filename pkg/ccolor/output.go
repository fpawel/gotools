package ccolor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type Output struct {
	f  *os.File
}

func NewWriter(f *os.File) io.Writer {
	return &Output{f: f}
}

func (x *Output) Write(p []byte) (int, error) {
	for _, p := range bytes.Split(p, []byte{'\n'}) {
		if len(p) == 0 {
			continue
		}
		Foreground(Green, true)
		if _, err := fmt.Fprint(x.f, time.Now().Format("15:04:05"), " "); err != nil {
			return 0, err
		}
		fields := bytes.Fields(p)
		if len(fields) > 1 {
			switch string(fields[1]) {
			case "ERR":
				Foreground(Red, true)
			case "WRN":
				Foreground(Yellow, true)
			case "inf":
				Foreground(White, true)
			default:
				Foreground(White, false)
			}
		}
		if _, err := x.f.Write(p); err != nil {
			return 0, err
		}
		ResetColor()
		if !bytes.HasSuffix(p, []byte("\n")) {
			if _, err := x.f.WriteString("\n"); err != nil {
				return 0, err
			}
		}
	}
	return len(p), nil
}
