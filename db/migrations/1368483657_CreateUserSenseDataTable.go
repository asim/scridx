package main

import (
    "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"

)

type UserSenseData struct {
  Id hood.Id
  UserId int64
  Password string `sql:"size(255),notnull"`
  Algo string `sql:"size(128),notnull"`
  Updated uint32
}

func (m *M) CreateUserSenseDataTable_1368483657_Up(hd *hood.Hood) {
  hd.CreateTable(&UserSenseData{})
}

func (m *M) CreateUserSenseDataTable_1368483657_Down(hd *hood.Hood) {
  hd.DropTable(&UserSenseData{})
}
