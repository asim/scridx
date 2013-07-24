package main

import (
	_ "github.com/ziutek/mymysql/godrv"
	"github.com/eaigner/hood"
)

func (m *M) AddCommentsIndex_1374154256_Up(hd *hood.Hood) {
	hd.CreateIndex("comments", "sdb_comments_entity_index", false, "entity_type", "entity_id")
	hd.CreateIndex("comments", "sdb_comments_karma_index", false, "karma")
}

func (m *M) AddCommentsIndex_1374154256_Down(hd *hood.Hood) {
	hd.DropIndex("comments", "sdb_comments_entity_index")
	hd.DropIndex("comments", "sdb_comments_karma_index")
}
