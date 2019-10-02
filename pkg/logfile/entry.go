package logfile

import (
	"errors"
	"time"
)

type Entry struct {
	Time time.Time `db:"tm"`
	Line string    `db:"tx"`
}

func parseEntry(line string, ent *Entry) error {
	if len(line) == 0 {
		return errors.New("empty line")
	}
	for i := range line {
		if line[i] == ' ' {
			if i+1 == len(line) {
				return errors.New("wrong format")
			}
			if t, err := time.Parse(layoutDatetime, line[:i]); err != nil {
				return err
			} else {
				ent.Line = line[i+1:]
				ent.Time = t
				return nil
			}
		}
	}
	return errors.New("wrong format")
}
