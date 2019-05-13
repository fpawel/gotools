package gorunex

import (
	"fmt"
	"github.com/fpawel/elco/pkg/winapp"
	"log"
	"os"
	"path/filepath"
	"time"
)

func LogFileName() string{
	exeDir := filepath.Dir(os.Args[0])
	t := time.Now()
	logDir := filepath.Join(exeDir, "logs")

	if err := winapp.EnsuredDirectory(logDir); err != nil {
		log.Fatal(err)
	}

	return filepath.Join(logDir, fmt.Sprintf("%s.log", t.Format("2006-01-02")))
}
