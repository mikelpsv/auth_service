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

func Install() {
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

	sql = `INSERT INTO clients (name, secret, key, expires, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = Db.Exec(sql, "Основная система", "hjhgHJGjhh767Kjh7", "1234567890", 3600, time.Now(), time.Now().Add(3600))
	if err != nil {
		log.Fatal(err.Error())
	}

	sql = `CREATE TABLE IF NOT EXISTS users	(
    			_id INTEGER PRIMARY KEY AUTOINCREMENT,
    			username  TEXT    default '',
    			password  TEXT    default '',
    			client_id INTEGER default 0
	);
`
	_, err = Db.Exec(sql)
	if err != nil {
		log.Fatal(err.Error())
	}

}
