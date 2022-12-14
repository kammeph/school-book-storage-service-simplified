package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbdriver   = os.Getenv("DB_DRIVER")
	dbuser     = os.Getenv("DB_USER")
	dbpassword = os.Getenv("DB_PASSWORD")
	dbprotocol = os.Getenv("DB_PROTOCOL")
	dbinstance = os.Getenv("DB_INSTANCE")
	dbdbname   = os.Getenv("DB_DATABASE")
)

func NewSqlDB() *sql.DB {
	// connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbuser, dbpassword, dbhost, dbport, dbdbname)
	connStr := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true", dbuser, dbpassword, dbprotocol, dbinstance, dbdbname)
	db, err := sql.Open(dbdriver, connStr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	log.Println("Successfully connected and pinged to sql database.")
	return db
}
