package main

import (
	"fmt"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gotools/pkg/logfile"
	"github.com/powerman/structlog"
	"os"
	"path/filepath"
)

func main() {

	log := structlog.New()

	structlog.DefaultLogger.
		SetPrefixKeys(
			structlog.KeyApp,
			structlog.KeyPID, structlog.KeyLevel, structlog.KeyUnit, structlog.KeyTime,
		).
		SetDefaultKeyvals(
			structlog.KeyApp, filepath.Base(os.Args[0]),
			structlog.KeySource, structlog.Auto,
		).
		SetSuffixKeys(
			structlog.KeyStack,
		).
		SetSuffixKeys(structlog.KeySource).
		SetKeysFormat(map[string]string{
			structlog.KeyTime:   " %[2]s",
			structlog.KeySource: " %6[2]s",
			structlog.KeyUnit:   " %6[2]s",
		})
	modbus.SetLogKeysFormat()

	if len(os.Args) < 2 {
		log.Fatalln("usage: run [exe name] [... exe args]")
	}
	name := os.Args[1]
	var arg []string
	if len(os.Args) > 2 {
		arg = os.Args[2:]
	}
	log.Info(fmt.Sprintf("%s %+v", name, arg))

	log.ErrIfFail(func() error {
		return logfile.Exec(nullWriter{}, name, arg...)
	})
}

type nullWriter struct{}

func (_ nullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
