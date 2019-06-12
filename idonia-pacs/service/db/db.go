package db

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

var globalDB *sql.DB

func NewDatabase(dbfile string) (database *sql.DB, err error) {
	database, err = sql.Open("sqlite3", dbfile)
	if err != nil {
		return
	}
	database.SetMaxOpenConns(1)
	var version string
	row := database.QueryRow("SELECT value FROM info WHERE key = 'version'")
	err = row.Scan(&version)
	if err != nil {
		err = nil
		for _, v := range INIT_DB {
			_, err = database.Exec(v)
			if err != nil {
				return database, err
			}
		}
	}
	for {
		var version string
		row := database.QueryRow("SELECT value FROM info WHERE key = 'version'")
		err := row.Scan(&version)
		if err != nil {
			break
		}
		queries, ok := MIGRATIONS_DB[version]
		if !ok {
			break
		}
		for _, v := range queries {
			_, _ = database.Exec(v)
		}
	}
	globalDB = database
	return globalDB, nil
}

func GetDB() (*sql.DB, error) {
	if globalDB != nil {
		return globalDB, nil
	}
	return nil, errors.New("No DB created")
}
