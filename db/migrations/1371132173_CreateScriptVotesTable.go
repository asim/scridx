package main

import (
  _ "github.com/ziutek/mymysql/godrv"
  "github.com/eaigner/hood"
)

type ScriptVotes struct {
  Id       hood.Id
  ScriptId int64 `sql:"notnull"`
  UserId   int64 `sql:"notnull"`
}

func (m *M) CreateScriptVotesTable_1371132173_Up(hd *hood.Hood) {
  hd.CreateTable(&ScriptVotes{})
}

func (m *M) CreateScriptVotesTable_1371132173_Down(hd *hood.Hood) {
  hd.DropTable(&ScriptVotes{})
}
