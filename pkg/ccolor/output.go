package ccolor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type Output struct {
	file *os.File
}

func NewWriter(file *os.File) io.Writer {
	return Output{file}
}

func (x Output) Write(p []byte) (int, error) {

	Foreground(Green, true)
	_, _ = fmt.Fprint(x.file, time.Now().Format("15:04:05"), " ")

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
	defer ResetColor()
	return x.file.Write(p)
}
