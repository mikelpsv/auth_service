package app

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var Db *sql.DB

func InitDb(host, dbname, dbuser, dbpass string) {
	var err error

	Db, err = sql.Open("sqlite3", ".data/authdata.db")
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	ctx := context.Background()
	err = Db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Close() {
	Db.Close()
}

func Install(loadTestData bool) {

	sql := `
		CREATE TABLE IF NOT EXISTS clients (
			_id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name TEXT,
			secret TEXT,
			key TEXT,
			expires INTEGER,
			created_at DATETIME,
			updated_at DATETIME
		)`
	_, err := Db.Exec(sql)
	if err != nil {
		log.Fatal(err.Error())
	}

	if loadTestData {
		sql = `INSERT INTO clients (name, secret, key, expires, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = Db.Exec(sql, "Основная система", "$2a$10$SyaL6fNLoPplhxqOlmN7MuA/MxXm7/F9AX.NqVDRSb4xi9YrHQg36", "1234567890", 3600, time.Now(), time.Now())
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	sql = `CREATE TABLE IF NOT EXISTS users	(
    			_id INTEGER PRIMARY KEY AUTOINCREMENT,
    			username  TEXT    default '',
    			password  TEXT    default '',
    			client_id INTEGER default 0,
    			created_at DATETIME,
				updated_at DATETIME
	);
`

	_, err = Db.Exec(sql)
	if err != nil {
		log.Fatal(err.Error())
	}

	if loadTestData {
		sql = `INSERT INTO users (username, password, client_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
		_, err = Db.Exec(sql, "Пользователь1", "$2a$10$/ui7v1gRNVLSRtfHOib/muwP5TAr7e33c9y7LPpfdUHmCIWJSO8ny", "1", time.Now(), time.Now())
		if err != nil {
			log.Fatal(err.Error())
		}

		sql = `INSERT INTO users (username, password, client_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
		_, err = Db.Exec(sql, "Пользователь2", "$2a$10$B2pAjD62tq0QOAswYaXqFe9cxVEgMm8PVTL4SfgIl3CNJUkmNITQm", "1", time.Now(), time.Now())
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
