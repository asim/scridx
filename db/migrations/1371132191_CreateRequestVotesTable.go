package main

import (
  _ "github.com/ziutek/mymysql/godrv"
  "github.com/eaigner/hood"
)

type RequestVotes struct {
  Id       hood.Id
  RequestId int64 `sql:"notnull"`
  UserId   int64 `sql:"notnull"`
}

func (m *M) CreateRequestVotesTable_1371132191_Up(hd *hood.Hood) {
  hd.CreateTable(&RequestVotes{})
}

func (m *M) CreateRequestVotesTable_1371132191_Down(hd *hood.Hood) {
  hd.DropTable(&RequestVotes{})
}
