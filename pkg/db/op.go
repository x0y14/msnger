package db

import (
	"github.com/x0y14/msnger/pkg/misc"
	"github.com/x0y14/msnger/pkg/protobuf"
)

type StoreOpReq struct {
	OpType protobuf.OperationType
	Param1 string
	Param2 string
	Param3 string
	Msg    *protobuf.Message
}

// StoreOp Operationの種類に従って、OpRelationをDBに登録していく。
func StoreOp(req *StoreOpReq) error {
	PingAndReconnect()

	switch req.OpType {
	case protobuf.OperationType_NOOP,

		protobuf.OperationType_CREATE_ACCOUNT,
		protobuf.OperationType_LOGIN,
		protobuf.OperationType_GET_EMAIL,

		protobuf.OperationType_CREATE_USER,
		protobuf.OperationType_GET_USER,

		protobuf.OperationType_FETCH_OPS,
		protobuf.OperationType_SEND_OP:
		// OperationType_NOOPはデバッグようなので実際使用できない。
		// revisionId発行
		rev := misc.GenerateRevisionId()
		// OPをDBに挿入
		err := insertOp(rev, req)
		if err != nil {
			return err
		}
		// OPが誰宛なのかDBに挿入
		err = insertOpRelation(req.Param1, rev)
		if err != nil {
			return err
		}
		// 宛先アカウントの最新revisionIdを更新
		err = updateLastRevision(req.Param1, rev)
		if err != nil {
			return err
		}
	case protobuf.OperationType_SEND_MESSAGE_SEND:
		// 送信したよOP用
		revSend := misc.GenerateRevisionId()
		// 受け取ったよOP用
		revRecv := misc.GenerateRevisionId()
		// 送信者のId
		sender := req.Param1
		// 受信者のId
		receiver := req.Param2
		// Param3はメッセージIDの予定

		// SendのOPをDBに保存入れちゃう
		err := insertOp(revSend, req)
		if err != nil {
			return err
		}
		// SendのOpは、送信者のみに送信
		err = insertOpRelation(sender, revSend)
		if err != nil {
			return err
		}
		// 送信者の最終OP更新
		err = updateLastRevision(sender, revSend)
		if err != nil {
			return err
		}

		// reqのOpTypeを、OperationType_SEND_MESSAGE_SEND -> OperationType_SEND_MESSAGE_RECV
		// に勝手に書き換えるよ。
		req.OpType = protobuf.OperationType_SEND_MESSAGE_RECV

		// RecvのOPをDBに保存
		err = insertOp(revRecv, req)
		if err != nil {
			return err
		}

		// 送信者と受信者に送信
		err = insertOpRelation(sender, revRecv)
		if err != nil {
			return err
		}
		err = insertOpRelation(receiver, revRecv)
		if err != nil {
			return err
		}

		// 送信者、受信者の最終OPを更新するよ。
		err = updateLastRevision(sender, revRecv)
		if err != nil {
			return err
		}
		err = updateLastRevision(receiver, revRecv)
		if err != nil {
			return err
		}

	case protobuf.OperationType_SEND_READ_RECEIPT_SEND:
		revSend := misc.GenerateRevisionId()
		revRecv := misc.GenerateRevisionId()
		sender := req.Param1
		receiver := req.Param2
		// readReceiptTarget := req.param3

		// 送信者が送信したというOPを保存
		err := insertOp(revSend, req)
		if err != nil {
			return err
		}
		// 通知
		err = insertOpRelation(sender, revSend)
		if err != nil {
			return err
		}
		// 最終OPを更新
		err = updateLastRevision(sender, revSend)
		if err != nil {
			return err
		}

		// reqのOpTypeを書き換えて、送信したという記録から、受信したという記録に書き換える。
		req.OpType = protobuf.OperationType_SEND_READ_RECEIPT_RECV

		// 受け取ったというOPを保存
		err = insertOp(revRecv, req)
		if err != nil {
			return err
		}

		// 送信者、受信者に通知
		err = insertOpRelation(sender, revRecv)
		if err != nil {
			return err
		}
		err = insertOpRelation(receiver, revRecv)
		if err != nil {
			return err

		}

		// 送信者、受信者の最終OPを更新
		err = updateLastRevision(sender, revRecv)
		if err != nil {
			return err
		}
		err = updateLastRevision(receiver, revRecv)
		if err != nil {
			return err
		}

	}
	return nil
}

func insertOp(revisionId uint64, req *StoreOpReq) error {
	ins, err := MsngerDB.Prepare("INSERT INTO msnger.Operation (revisionId, type, param1, param2, param3, messageId) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer ins.Close()
	msgId := ""
	if req.Msg != nil {
		msgId = req.Msg.Id
	}
	_, err = ins.Exec(revisionId, req.OpType, req.Param1, req.Param2, req.Param3, msgId)
	return err
}

func insertOpRelation(userId string, revisionId uint64) error {
	ins, err := MsngerDB.Prepare("INSERT INTO msnger.OpRelation (targetUserId, revisionId) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(userId, revisionId)
	return err
}

func updateLastRevision(userId string, lastRevision uint64) error {
	ins, err := MsngerDB.Prepare("INSERT INTO msnger.OpRelation (targetUserId, revisionId) values (?, ?) on duplicate key update revisionId = ?")
	if err != nil {
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(userId, lastRevision, lastRevision)
	return err
}
