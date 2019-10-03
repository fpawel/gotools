package logfile

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func newDB() *sqlx.DB {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(0)
	db := sqlx.NewDb(conn, "sqlite3")
	db.MustExec(`
CREATE TABLE entry(
	tm DATETIME  NOT NULL ,
	tx TEXT NOT NULL    
)`)

	return db
}
func dbInsertEntry(db *sqlx.DB, ent Entry) {
	db.MustExec(`INSERT INTO entry(tm, tx) VALUES (?,?)`, ent.Time, ent.Line)
}

func dbSelectEntries(db *sqlx.DB, filter string, entries *[]Entry) error {
	return db.Select(&entries, "SELECT * FROM entry"+filter+" ORDER BY tm")
}
