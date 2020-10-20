package app

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

var Db *sql.DB

func InitDb(host, dbname, dbuser, dbpass string) {
	var err error
/*
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, dbuser, dbpass)
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

*/

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;", host, dbuser, dbpass, dbname)

	Db, err = sql.Open("sqlserver", connString)
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
