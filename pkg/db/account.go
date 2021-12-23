package db

import (
	"database/sql"
	"github.com/x0y14/msnger/pkg/auth"
	"github.com/x0y14/msnger/pkg/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type InsertAccountReq struct {
	Id       string
	Email    string
	Password string
	isAdmin  bool
}

type SelectAccountReq struct {
	Id    string
	Email string
}

type AccountData struct {
	Id        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	isAdmin   []byte    `db:"isAdmin"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}

func InsertAccount(req *InsertAccountReq) (string, error) {
	// 接続確認
	PingAndReconnect()

	// jwt token作成
	token, err := auth.GenerateJWTToken(req.Id)
	if err != nil {
		return "", err
	}

	// データ挿入
	ins, err := MsngerDB.Prepare("INSERT INTO msnger.Account (id, email, password, isAdmin) VALUES (?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	_, err = ins.Exec(req.Id, req.Email, req.Password, req.isAdmin)
	if err != nil {
		return "", err
	}

	return token, nil
}

func SelectAccountWithId(req *SelectAccountReq) (*protobuf.Account, error) {
	row := MsngerDB.QueryRow(`select id, email, password, isAdmin, createdAt, updatedAt from msnger.Account where id = ? limit 1`, req.Id)

	if row.Err() == sql.ErrNoRows {
		return nil, nil
	}

	if row.Err() != nil && row.Err() != sql.ErrNoRows {
		return nil, row.Err()
	}

	ad := AccountData{}
	if err := row.Scan(&ad.Id, &ad.Email, &ad.Password, &ad.isAdmin, &ad.CreatedAt, &ad.UpdatedAt); err != nil {
		return nil, err
	}

	var isAdmin bool
	if ad.isAdmin[0] == 0 {
		isAdmin = false
	} else {
		isAdmin = true
	}

	account := protobuf.Account{
		Id:        ad.Id,
		Email:     ad.Email,
		Password:  ad.Password,
		IsAdmin:   isAdmin,
		CreatedAt: timestamppb.New(ad.CreatedAt),
		UpdatedAt: timestamppb.New(ad.UpdatedAt),
	}

	return &account, nil
}
func SelectAccountWithEmail(req *SelectAccountReq) (*protobuf.Account, error) {
	panic("unimplemented")
}
