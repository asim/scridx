package main

import (
    "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)


func (m *M) AddCommentVotesIndex_1371134206_Up(hd *hood.Hood) {
  hd.CreateIndex("comment_votes", "sdb_comment_votes_comment_id_user_id_index", true, "comment_id", "user_id")
}

func (m *M) AddCommentVotesIndex_1371134206_Down(hd *hood.Hood) {
  hd.DropIndex("comment_votes", "sdb_comment_votes_comment_id_user_id_index")
}
