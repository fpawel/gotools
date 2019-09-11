package data

import (
	"bytes"
	"github.com/fpawel/gohelp"
	"github.com/jmoiron/sqlx"
	"github.com/powerman/structlog"
	"os"
	"path/filepath"
	"sync"
)

//go:generate go run github.com/fpawel/gotools/cmd/sqlstr/...

type DB struct {
	db      *sqlx.DB
	kind    Kind
	mu      sync.Mutex
	exeName string
	s       string
}

type Kind int

const (
	KindText Kind = iota
	KindJSON
)

func New(kind Kind, exeName string) *DB {
	x := &DB{db: gohelp.OpenSqliteDBx(filepath.Join("loggo.sqlite")), kind: kind, exeName: exeName}
	x.db.MustExec(SQLCreate)
	return x
}

func (x *DB) Close() error {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.db.Close()
}

func (x *DB) WriteRawMessage(msg string) {
	x.mu.Lock()
	defer x.mu.Unlock()
	r := x.db.MustExec(`INSERT INTO entry(msg) VALUES (?)`, msg)
	entryID, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}
	x.db.MustExec(`INSERT INTO meta(entry_id, tag, value) VALUES (?,?,?), (?,?,?)`,
		entryID, structlog.KeyApp, x.exeName,
		entryID, structlog.KeyPID, os.Getpid(),
	)
}

func (x *DB) Write(p []byte) (int, error) {
	if !bytes.HasSuffix(p, []byte("\n")) {
		x.s += string(p)
		return len(p), nil
	}
	msg := x.s + string(p)
	x.s = ""
	go x.WriteRawMessage(msg)

	return len(p), nil
}
