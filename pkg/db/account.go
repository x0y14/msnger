package db

import (
	"database/sql"
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

// InsertAccount Return JwtToken, error
func InsertAccount(req *InsertAccountReq) error {
	// 接続確認
	PingAndReconnect()

	// データ挿入
	ins, err := MsngerDB.Prepare("INSERT INTO msnger.Account (id, email, password, isAdmin) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer ins.Close()
	_, err = ins.Exec(req.Id, req.Email, req.Password, req.isAdmin)
	if err != nil {
		return err
	}

	return nil
}

func SelectAccountWithId(req *SelectAccountReq) (*protobuf.Account, error) {
	PingAndReconnect()

	row := MsngerDB.QueryRow(`select id, email, password, isAdmin, createdAt, updatedAt from msnger.Account where id = ? limit 1`, req.Id)

	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}
		return nil, row.Err()
	}

	ad := AccountData{}
	if err := row.Scan(&ad.Id, &ad.Email, &ad.Password, &ad.isAdmin, &ad.CreatedAt, &ad.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
	PingAndReconnect()

	row := MsngerDB.QueryRow(`select id, email, password, isAdmin, createdAt, updatedAt from msnger.Account where email = ? limit 1`, req.Email)

	if row.Err() != nil {
		if row.Err() == sql.ErrNoRows {
			return nil, nil
		}
		return nil, row.Err()
	}

	ad := AccountData{}
	if err := row.Scan(&ad.Id, &ad.Email, &ad.Password, &ad.isAdmin, &ad.CreatedAt, &ad.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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
