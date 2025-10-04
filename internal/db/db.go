package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"queueit/internal/helper"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed migration/schema.sql
var schema string

const SQLLITE_DB_FILE_NAME = "queueit.db"

var dbinfo *DBInfo

type DBInfo struct {
	conn   *sql.DB
	dbfile string
}

func InitDB() error {

	// pre-checks for db file
	loc, err := helper.GetAppDataDir()
	if err != nil {
		return err
	}

	if _, err = os.Stat(loc); os.IsNotExist(err) {
		if err = os.MkdirAll(loc, 0755); err != nil {
			return err
		}
	}

	loc = filepath.Join(loc, SQLLITE_DB_FILE_NAME)
	if _, err = os.Stat(loc); os.IsNotExist(err) {
		if err = os.WriteFile(loc, []byte(""), 0755); err != nil {
			return err
		}
	}

	// init db connection
	conn, err := sql.Open("sqlite", loc)
	if err != nil {
		return err
	}

	// check and create sql-schema
	if _, err := conn.Exec(string(schema)); err != nil {
		return err
	}

	setDBInfo(&DBInfo{
		conn:   conn,
		dbfile: loc,
	})
	return nil
}

func (di *DBInfo) E(query string, args ...any) (sql.Result, error) {
	return di.conn.Exec(query, args...)
}

func (di *DBInfo) Q(query string, args ...any) (*sql.Rows, error) {
	return di.conn.Query(query, args...)
}

func GetDBInfo() *DBInfo   { return dbinfo }
func setDBInfo(di *DBInfo) { dbinfo = di }
