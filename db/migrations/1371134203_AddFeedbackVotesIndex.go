package main

import (
    "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddFeedbackVotesIndex_1371134203_Up(hd *hood.Hood) {
  hd.CreateIndex("feedback_votes", "sdb_feedback_votes_feedback_id_user_id_index", true, "feedback_id", "user_id")
}

func (m *M) AddFeedbackVotesIndex_1371134203_Down(hd *hood.Hood) {
  hd.DropIndex("feedback_votes", "sdb_feedback_votes_feedback_id_user_id_index")
}
