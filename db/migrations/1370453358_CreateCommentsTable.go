package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

type Comments struct {
  Id hood.Id
  Text string `sql:"size(65535),notnull"`
  Rank float64 `sql:"default(0.0),notnull"`
  Karma int32 `sql:"default(1),notnull"`
  Stored uint32 `sql:"notnull"`
  UserId int64 `sql:"notnull"`
  ParentId int64 `sql:"notnull"`
  EntityId int64 `sql:"notnull"`
  EntityType int16 `sql:"notnull"`
}

func (m *M) CreateCommentsTable_1370453358_Up(hd *hood.Hood) {
  hd.CreateTable(&Comments{})
}

func (m *M) CreateCommentsTable_1370453358_Down(hd *hood.Hood) {
  hd.DropTable(&Comments{})
}
