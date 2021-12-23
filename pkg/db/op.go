package db

import (
	"github.com/x0y14/msnger/pkg/misc"
	"github.com/x0y14/msnger/pkg/protobuf"
	"log"
)

// StoreOp Operationの種類に従って、必要とする人に向けてOperationInfoを配る。OperationInfoは名前を変えた方がいいと思う。
func StoreOp(op *protobuf.Operation) {
	// connection check
	PingAndReconnect()

	switch op.Type {
	case protobuf.OperationType_NOOP:
		revisionId := misc.GenerateRevisionId()
		// Operation
		OpInsert, err := MsngerDB.Prepare("INSERT INTO msnger.Operation (revisionId, type, param1, param2, param3, msgId) values (?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatalln(err)
		}
		defer OpInsert.Close()
		_, err = OpInsert.Exec(revisionId, op.GetType(), op.GetParam1(), op.GetParam2(), op.GetParam3(), "")
		if err != nil {
			log.Fatalln(err)
		}

		// OperationInfo
		InfoInsert, err := MsngerDB.Prepare("INSERT into msnger.OperationInfo (id, revisionId, target) values (?, ?, ?)")
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
