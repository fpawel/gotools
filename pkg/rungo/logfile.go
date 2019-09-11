package rungo

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func newLogFileOutput() io.WriteCloser {
	logFile, err := os.OpenFile(LogFileName(), os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return &logFileOutput{f: logFile}
}

type logFileOutput struct {
	f  *os.File
	mu sync.Mutex
}

func (x *logFileOutput) Close() error {
	return x.f.Close()
}

func (x *logFileOutput) Write(p []byte) (int, error) {
	go func() {
		x.mu.Lock()
		defer x.mu.Unlock()
		if _, err := fmt.Fprint(x.f, time.Now().Format("15:04:05"), " "); err != nil {
			log.PrintErr(err)
		}
		if _, err := x.f.Write(p); err != nil {
			log.PrintErr(err)
		}
	}()
	return len(p), nil
}
