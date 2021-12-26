package db

import (
	"fmt"
	"github.com/x0y14/msnger/pkg/misc"
	"github.com/x0y14/msnger/pkg/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

//type StoreOpReq struct {
//	OpType protobuf.OperationType
//	Param1 string
//	Param2 string
//	Param3 string
//	Msg    *protobuf.Message
//}

type OperationData struct {
	RevisionId uint64    `db:"revisionId"`
	Type       int       `db:"type"`
	Param1     string    `db:"param1"`
	Param2     string    `db:"param2"`
	Param3     string    `db:"param3"`
	MessageId  string    `db:"messageId"`
	CreatedAt  time.Time `db:"createdAt"`
	UpdatedAt  time.Time `db:"updatedAt"`
}

// StoreOp Operationの種類に従って、OpRelationをDBに登録していく。
func StoreOp(req *protobuf.Operation) error {
	PingAndReconnect()

	switch req.Type {
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
		req.Type = protobuf.OperationType_SEND_MESSAGE_RECV

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
		req.Type = protobuf.OperationType_SEND_READ_RECEIPT_RECV

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

func insertOp(revisionId uint64, req *protobuf.Operation) error {
	ins, err := MsngerDB.Prepare("INSERT INTO msnger.Operation (revisionId, type, param1, param2, param3, messageId) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer ins.Close()
	var msgId string
	if req.Message == nil {
		msgId = ""
	} else {
		msgId = req.Param3
	}
	_, err = ins.Exec(revisionId, req.Type, req.Param1, req.Param2, req.Param3, msgId)
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
	ins, err := MsngerDB.Prepare("INSERT INTO msnger.LastRevision (id, lastRevisionId) values (?, ?) on duplicate key update lastRevisionId = ?")
	if err != nil {
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(userId, lastRevision, lastRevision)
	return err
}

func GetLastOpRevision(userId string) (uint64, error) {
	PingAndReconnect()

	row := MsngerDB.QueryRow(`select lastRevisionId from msnger.LastRevision where id = ? limit 1`, userId)
	var revisionId uint64
	err := row.Scan(&revisionId)
	if err != nil {
		return 0, fmt.Errorf("failed to get revisionId of %v: %v", userId, err)
	}
	return revisionId, nil
}

//func GetOp(revisionId uint64) {}

func GetOpsBiggerThan(userId string, revisionId uint64) ([]*protobuf.Operation, error) {
	PingAndReconnect()

	// todo : Opで検索するのではなく、OpRelationで検索して、見つけたIDでOPを探す。
	rows, err := MsngerDB.Query(`select revisionId, type, param1, param2, param3, messageId, createdAt, updatedAt from msnger.Operation where (param1 = ? or param2 = ?) and ? < revisionId`, userId, userId, revisionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ops []*protobuf.Operation

	for rows.Next() {
		// 一時データに変換
		var opData OperationData
		err = rows.Scan(&opData.RevisionId, &opData.Type, &opData.Param1, &opData.Param2, &opData.Param3, &opData.MessageId, &opData.CreatedAt, &opData.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to get op %v: %v", userId, err)
		}
		// protobufのデータに変換
		var msg *protobuf.Message
		if opData.MessageId != "" {
			msg, err = GetMessage(opData.MessageId)
			if err != nil {
				return nil, err
			}
		}
		log.Printf("%v", msg.String())
		op := protobuf.Operation{
			RevisionId: opData.RevisionId,
			Type:       protobuf.OperationType(opData.Type),
			Param1:     opData.Param1,
			Param2:     opData.Param2,
			Param3:     opData.Param3,
			Message:    msg,
			CreatedAt:  timestamppb.New(opData.CreatedAt),
			UpdatedAt:  timestamppb.New(opData.UpdatedAt),
		}
		ops = append(ops, &op)
	}

	return ops, nil
}
