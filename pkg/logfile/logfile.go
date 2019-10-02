package logfile

import (
	"bufio"
	"fmt"
	"github.com/fpawel/gohelp/must"
	"github.com/jmoiron/sqlx"
	"github.com/powerman/structlog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func ListDays() []time.Time {
	r := regexp.MustCompile(`\d\d\d\d-\d\d-\d\d`)
	m := make(map[time.Time]struct{})
	_ = filepath.Walk(logDir, func(path string, f os.FileInfo, _ error) error {
		if f == nil || f.IsDir() {
			return nil
		}
		s := r.FindString(f.Name())
		if len(s) == 0 {
			return nil
		}
		t, err := time.Parse(layoutDate, s)
		if err != nil {
			return nil
		}
		m[daytime(t)] = struct{}{}
		return nil
	})
	var days []time.Time
	for t := range m {
		days = append(days, t)
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Before(days[j])
	})
	return days
}

func Read(t time.Time, filter string) ([]Entry, error) {
	t = daytime(t)
	re := regexp.MustCompile(`\d\d\d\d-\d\d-\d\d`)
	db := newDB()
	defer log.ErrIfFail(db.Close)

	filter = strings.TrimSpace(filter)
	if len(filter) != 0 {
		filter = " WHERE " + filter
	}
	if err := dbSelectEntries(db, filter, &[]Entry{}); err != nil {
		return nil, err
	}
	_ = filepath.Walk(logDir, func(path string, f os.FileInfo, _ error) error {
		if f == nil || f.IsDir() {
			return nil
		}
		strTime := re.FindString(f.Name())
		if len(strTime) == 0 {
			return nil
		}
		if fileTime, err := time.Parse(layoutDate, strTime); err == nil && fileTime == t {
			readEntries(path, db)
		}
		return nil
	})

	var entries []Entry
	if err := dbSelectEntries(db, filter, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func readEntries(filename string, db *sqlx.DB) {
	file, err := os.Open(filename)

	if os.IsNotExist(err) {
		return
	}

	if err != nil {
		log.PrintErr(err, "file", filepath.Base(filename))
		return
	}
	defer log.ErrIfFail(file.Close)

	var (
		lineNumber int
		scanner    = bufio.NewScanner(file)
	)

	for scanner.Scan() {
		line := scanner.Text()
		var ent Entry
		if err := parseEntry(line, &ent); err != nil {
			log.PrintErr(err,
				"line", fmt.Sprintf("%d:`%s`", lineNumber, line),
				"file", filepath.Base(filename))
			continue
		}
		dbInsertEntry(db, ent)
		lineNumber++
	}
	must.AbortIf(scanner.Err())
}

func filename(t time.Time, suffix string) string {
	return filepath.Join(logDir, fmt.Sprintf("%s%s.log", t.Format(layoutDate), suffix))
}

func daytime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func ensureDir() error {
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) { // создать каталог если его нет
		err = os.MkdirAll(logDir, os.ModePerm)
	}
	return err
}

var (
	log    = structlog.New()
	logDir = filepath.Join(filepath.Dir(os.Args[0]), "logs")
)

const (
	layoutDatetime = "2006-01-02-15:04:05.000"
	layoutDate     = "2006-01-02"
)
