package main

import (
	"database/sql"
	"log"
)

var wamdb *sql.DB

func openWAM() {
	var err error
	wamdb, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
		CREATE TABLE wam ( MDate TEXT, MTime TEXT, MMedia TEXT, MSender TEXT)
		`
	_, err = wamdb.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

}

func insertWAM(mtype string, mdate string, mtime string, mmedia string, msender string) {
	stmt, err := wamdb.Prepare("insert into wam(MDate, MTime, MMedia, MSender) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(mdate, mtime, mmedia, msender)
	if err != nil {
		log.Fatal(err)
	}

}
