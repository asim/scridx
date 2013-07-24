package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddRequestIndex_1374154257_Up(hd *hood.Hood) {
  hd.CreateIndex("requests", "sdb_requests_title_index", false, "title")
  hd.CreateIndex("requests", "sdb_requests_stored_index", false, "stored")
  hd.CreateIndex("requests", "sdb_requests_drafted_index", false, "drafted")
  hd.CreateIndex("requests", "sdb_requests_source_index", true, "source")
  hd.CreateIndex("requests", "sdb_requests_rank_index", false, "rank")
}

func (m *M) AddRequestIndex_1374154257_Down(hd *hood.Hood) {
  hd.DropIndex("requests", "sdb_requests_title_index")
  hd.DropIndex("requests", "sdb_requests_stored_index")
  hd.DropIndex("requests", "sdb_requests_drafted_index")
  hd.DropIndex("requests", "sdb_requests_source_index")
  hd.DropIndex("requests", "sdb_requests_rank_index")
}
