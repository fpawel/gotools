package logfile

import (
	"github.com/fpawel/gotools/pkg/ccolor"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func Exec(addWriter io.Writer, name string, arg ...string) error {
	logFileOutput, err := NewOutput("." + filepath.Base(name))
	if err != nil {
		return err
	}
	defer log.ErrIfFail(logFileOutput.Close)

	wrt := func(f *os.File) io.Writer {
		if addWriter == nil {
			return io.MultiWriter(logFileOutput, ccolor.NewWriter(f))
		}
		return io.MultiWriter(logFileOutput, ccolor.NewWriter(f), addWriter)
	}

	cmd := exec.Command(name, arg...)
	cmd.Stderr = wrt(os.Stderr)
	cmd.Stdout = wrt(os.Stdout)

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
