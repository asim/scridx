package main

import (
  "github.com/eaigner/hood"
  _ "github.com/ziutek/mymysql/godrv"
)

func (m *M) AddFeedbackIndex_1374154259_Up(hd *hood.Hood) {
  hd.CreateIndex("feedback", "sdb_feedback_title_index", false, "title")
  hd.CreateIndex("feedback", "sdb_feedback_source_index", true, "source")
  hd.CreateIndex("feedback", "sdb_feedback_rank_index", false, "rank")
}

func (m *M) AddFeedbackIndex_1374154259_Down(hd *hood.Hood) {
  hd.DropIndex("feedback", "sdb_feedback_title_index")
  hd.DropIndex("feedback", "sdb_feedback_source_index")
  hd.DropIndex("feedback", "sdb_feedback_rank_index")
}
