package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/x0y14/msnger/pkg/protobuf"
	op2 "github.com/x0y14/msnger/pkg/service/op"
	"log"
	"time"
)

var OperationDB *sql.DB

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "msnger_admin:password@/msnger")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func init() {
	OperationDB = NewDB()
}

func PingAndReconnect() {
	// try ping...
	if err := OperationDB.Ping(); err != nil {
		log.Printf("OperationDB->Ping Err: %v", err)
		// failed to connect...
	} else {
		// connection is alive!
		return
	}

	// try close connection
	err := OperationDB.Close()
	if err != nil {
		log.Printf("OperationDB->Close Err: %v", err)
		// failed to close connection
	}

	// force replace
	OperationDB = NewDB()
}

func StoreOp(op *protobuf.Operation) {
	// connection check
	PingAndReconnect()

	switch op.Type {
	case protobuf.OperationType_NOOP:
		revisionId := op2.GenerateRevisionId()
		// Operation
		OpInsert, err := OperationDB.Prepare("INSERT INTO msnger.Operation (revisionId, type, param1, param2, param3, msgId) values (?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatalln(err)
		}
		defer OpInsert.Close()
		_, err = OpInsert.Exec(revisionId, op.GetType(), op.GetParam1(), op.GetParam2(), op.GetParam3(), "")
		if err != nil {
			log.Fatalln(err)
		}

		// OperationInfo
		InfoInsert, err := OperationDB.Prepare("INSERT into msnger.OperationInfo (id, revisionId, target) values (?, ?, ?)")
		if err != nil {
			log.Fatalln(err)
		}
		defer InfoInsert.Close()
		_, err = InfoInsert.Exec(0, revisionId, op.GetParam1())
		if err != nil {
			log.Fatalln(err)
		}
	}
}
