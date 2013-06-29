package main

import (
    "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddRequestVotesIndex_1371134205_Up(hd *hood.Hood) {
  hd.CreateIndex("request_votes", "sdb_request_votes_request_id_user_id_index", true, "request_id", "user_id")
}

func (m *M) AddRequestVotesIndex_1371134205_Down(hd *hood.Hood) {
  hd.DropIndex("request_votes", "sdb_request_votes_request_id_user_id_index")
}
