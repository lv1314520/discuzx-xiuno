package database

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/skiy/xiuno-tools/app/libraries/mcfg"
)

func GetUcDB() (db gdb.DB) {
	db = g.DB("uc")
	db.SetDebug(mcfg.GetCfg().GetBool("database.uc.0.debug"))
	return db
}

func GetXiunoDB() (db gdb.DB) {
	db = g.DB("xiuno")
	db.SetDebug(mcfg.GetCfg().GetBool("database.xiuno.0.debug"))
	return db
}

func GetDiscuzDB() (db gdb.DB) {
	db = g.DB("discuz")
	db.SetDebug(mcfg.GetCfg().GetBool("database.discuz.0.debug"))
	return db
}
