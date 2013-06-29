package main

import (
    "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddScriptVotesIndex_1371134202_Up(hd *hood.Hood) {
  hd.CreateIndex("script_votes", "sdb_script_votes_script_id_user_id_index", true, "script_id", "user_id")
}

func (m *M) AddScriptVotesIndex_1371134202_Down(hd *hood.Hood) {
  hd.DropIndex("script_votes", "sdb_script_votes_script_id_user_id_index")
}
