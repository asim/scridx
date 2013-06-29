package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddScriptIndex_1367258614_Up(hd *hood.Hood) {
  hd.CreateIndex("scripts", "sdb_scripts_title_index", false, "title")
  hd.CreateIndex("scripts", "sdb_scripts_stored_index", false, "stored")
  hd.CreateIndex("scripts", "sdb_scripts_drafted_index", false, "drafted")
  hd.CreateIndex("scripts", "sdb_scripts_source_index", true, "source")
}

func (m *M) AddScriptIndex_1367258614_Down(hd *hood.Hood) {
  hd.DropIndex("scripts", "sdb_scripts_title_index")
  hd.DropIndex("scripts", "sdb_scripts_stored_index")
  hd.DropIndex("scripts", "sdb_scripts_drafted_index")
  hd.DropIndex("scripts", "sdb_scripts_source_index")
}
