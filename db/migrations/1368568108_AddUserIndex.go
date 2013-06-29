package main

import (
    "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddUserIndex_1368568108_Up(hd *hood.Hood) {
  hd.CreateIndex("users", "sdb_users_users_username_index", true, "username")
}

func (m *M) AddUserIndex_1368568108_Down(hd *hood.Hood) {
  hd.DropIndex("users", "sdb_users_users_username_index")
}
