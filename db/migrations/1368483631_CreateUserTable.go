package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

type Users struct {
  Id       hood.Id
  Username string `sql:"size(15),notnull"`
  Email    string `sql:"size(255),notnull"`
  Name     string `sql:"size(20),notnull"`
  Logline  string `sql:"size(160)"`
  Karma    int32 `sql:"default(0)"`
  Updated  uint32
  Created  uint32
  LoginAttempts int `sql:"default(0),notnull"`
  Locked   bool `sql:"default(0),notnull"`
}

func (m *M) CreateUserTable_1368483631_Up(hd *hood.Hood) {

  hd.CreateTable(&Users{})
}

func (m *M) CreateUserTable_1368483631_Down(hd *hood.Hood) {
  hd.DropTable(&Users{})
}
