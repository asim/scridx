package main

import (
    "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddUserSenseDataIndex_1368568433_Up(hd *hood.Hood) {
  hd.CreateIndex("user_sense_data", "sdb_users_user_sense_data_user_id_index", true, "user_id")
}

func (m *M) AddUserSenseDataIndex_1368568433_Down(hd *hood.Hood) {
  hd.DropIndex("user_sense_data", "sdb_users_user_sense_data_user_id_index")
}
