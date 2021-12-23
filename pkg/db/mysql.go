package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var MsngerDB *sql.DB

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "msnger_admin:password@/msnger?parseTime=true")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func init() {
	MsngerDB = NewDB()
}

func PingAndReconnect() {
	// try ping...
	if err := MsngerDB.Ping(); err != nil {
		log.Printf("MsngerDB->Ping Err: %v", err)
		// failed to connect...
	} else {
		// connection is alive!
		return
	}

	// try close connection
	err := MsngerDB.Close()
	if err != nil {
		log.Printf("MsngerDB->Close Err: %v", err)
		// failed to close connection
	}

	// force replace
	MsngerDB = NewDB()
}
