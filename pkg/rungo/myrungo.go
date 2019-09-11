package rungo

import (
	"bytes"
	"fmt"
	"github.com/fpawel/gohelp/winapp"
	"github.com/fpawel/gotools/pkg/ccolor"
	"github.com/maruel/panicparse/stack"
	"github.com/powerman/structlog"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func LogFileName() string {
	exeDir := filepath.Dir(os.Args[0])
	t := time.Now()
	logDir := filepath.Join(exeDir, "logs")

	if err := winapp.EnsuredDirectory(logDir); err != nil {
		log.Fatal(err)
	}

	return filepath.Join(logDir, fmt.Sprintf("%s.log", t.Format("2006-01-02")))
}

func Process(exeName string, args string, onPanic func(), writers ...io.Writer) error {
	logFile := newLogFileOutput()

	defer structlog.New().ErrIfFail(logFile.Close)

	panicOutput := bytes.NewBuffer(nil)

	writers = append(writers, logFile, panicOutput)

	cmd := exec.Command(exeName, strings.Fields(args)...)
	cmd.Stderr = io.MultiWriter(append(writers, ccolor.NewWriter(os.Stderr))...)
	cmd.Stdout = io.MultiWriter(append(writers, ccolor.NewWriter(os.Stdout))...)
	if err := cmd.Start(); err != nil {
		return err
	}
	err := cmd.Wait()
	if err == nil {
		return nil
	}
	if onPanic != nil {
		onPanic()
	}
	if _, err := fmt.Fprintln(cmd.Stderr, err); err != nil {
		return err
	}
	panicContent := bytes.NewBuffer(nil)
	if err := parseDump(panicOutput, panicContent); err != nil {
		return fmt.Errorf("unknown panic: %v", err)
	}
	if _, err := io.WriteString(cmd.Stderr, "panic occurred!\n"); err != nil {
		return err
	}
	if _, err := panicContent.WriteTo(cmd.Stderr); err != nil {
		return err
	}
	return nil
}

func parseDump(in io.Reader, out io.Writer) error {
	// Optional: Check for GOTRACEBACK being set, in particular if there is only
	// one goroutine returned.
	c, err := stack.ParseDump(in, out, true)
	if err != nil {
		return err
	}

	// Find out similar goroutine traces and group them into buckets.
	buckets := stack.Aggregate(c.Goroutines, stack.AnyValue)

	// Calculate alignment.
	srcLen := 0
	pkgLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(line.SrcLine()); l > srcLen {
				srcLen = l
			}
			if l := len(line.Func.PkgName()); l > pkgLen {
				pkgLen = l
			}
		}
	}

	for _, bucket := range buckets {
		// Print the goroutine header.
		extra := ""
		if s := bucket.SleepString(); s != "" {
			extra += " [" + s + "]"
		}
		if bucket.Locked {
			extra += " [locked]"
		}
		if c := bucket.CreatedByString(false); c != "" {
			extra += " [Created by " + c + "]"
		}
		if _, err := fmt.Fprintf(out, "%d: %s%s\n", len(bucket.IDs), bucket.State, extra); err != nil {
			return err
		}

		// Print the stack lines.
		for _, line := range bucket.Stack.Calls {
			if _, err := fmt.Fprintf(out,
				"    %-*s %-*s %s(%s)\n",
				pkgLen, line.Func.PkgName(), srcLen, line.SrcLine(),
				line.Func.Name(), &line.Args); err != nil {
				return err
			}
		}
		if bucket.Stack.Elided {
			if _, err := fmt.Fprintf(out, "    (...)\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func newLogFileOutput() io.WriteCloser {
	logFile, err := os.OpenFile(LogFileName(), os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return logFileOutput{logFile}
}

type logFileOutput struct {
	logFile *os.File
}

func (x logFileOutput) Close() error {
	return x.logFile.Close()
}

func (x logFileOutput) Write(p []byte) (int, error) {
	if _, err := fmt.Fprint(x.logFile, time.Now().Format("15:04:05"), " "); err != nil {
		return 0, err
	}
	return x.logFile.Write(p)
}
