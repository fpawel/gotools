package main

import (
	"bytes"
	"fmt"
	"github.com/fpawel/gorunex/pkg/gorunex"
	"io"
	"log"
	"os"
	"time"
)

type output struct {
	logFile     *os.File
	panicBuffer *bytes.Buffer
}

func NewOutput() output{
	logFile, err := os.OpenFile(gorunex.LogFileName(), os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return output{
		panicBuffer : bytes.NewBuffer(nil),
		logFile:logFile,
	}
}

func (x output) PrintPanic() error {
	panicParsed := bytes.NewBuffer(nil)
	if err := parseDump(x.panicBuffer, panicParsed); err != nil {
		return fmt.Errorf("unknown panic: %v", err)
	}
	if _,err := io.Copy(output{},panicParsed ); err != nil {
		return err
	}
	return nil
}

func (x output) Close() error {
	return x.logFile.Close()
}

func (x output) Write(p []byte) (int, error) {

	Foreground(Green, true)
	fmt.Print(time.Now().Format("15:04:05"), " ")

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
	_, _ = os.Stderr.Write(p)

	ResetColor()

	_, _ = fmt.Fprint(x.logFile, time.Now().Format("15:04:05"), " ")
	_,_ = x.panicBuffer.Write(p)
	return x.logFile.Write(p)
}
