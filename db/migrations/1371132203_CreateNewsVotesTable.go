package main

import (
  _ "github.com/ziutek/mymysql/godrv"
  "github.com/eaigner/hood"
)

type NewsVotes struct {
  Id       hood.Id
  NewsId int64 `sql:"notnull"`
  UserId   int64 `sql:"notnull"`
}

func (m *M) CreateNewsVotesTable_1371132203_Up(hd *hood.Hood) {
  hd.CreateTable(&NewsVotes{})
}

func (m *M) CreateNewsVotesTable_1371132203_Down(hd *hood.Hood) {
  hd.DropTable(&NewsVotes{})
}
