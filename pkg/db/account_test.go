package db_test

import (
	"fmt"
	"github.com/x0y14/msnger/pkg/db"
	"github.com/x0y14/msnger/pkg/misc"
	"log"
	"testing"
)

func TestInsertAccount(t *testing.T) {
	req := db.InsertAccountReq{
		Id:       misc.GenerateUserId(),
		Email:    "sample2@example.com",
		Password: "123456",
	}
	jwt, err := db.InsertAccount(&req)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("jwt: %v", jwt)
}

func TestSelectAccountWithEmail(t *testing.T) {
	req := db.SelectAccountReq{
		Id:    "",
		Email: "sample@example.com",
	}
	account, err := db.SelectAccountWithEmail(&req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%v", account.String())
}

func TestSelectAccountWithId(t *testing.T) {
	req := db.SelectAccountReq{
		Id:    "Uc725f4is1s45mtd6ic6g",
		Email: "sample@example.com",
	}
	account, err := db.SelectAccountWithId(&req)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Id: %v\n", account.Id)
	fmt.Printf("Email: %v\n", account.Email)
	fmt.Printf("Password: %v\n", account.Password)
	fmt.Printf("isAdmin: %v\n", account.IsAdmin)
	fmt.Printf("CreatedAt: %v\n", account.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", account.UpdatedAt)
	//log.Printf("%v", account.String())
}
