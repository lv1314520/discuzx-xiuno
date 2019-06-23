package controllers

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/util/gconv"
	"time"
	"xiuno-tools/app/libraries/common"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mlog"
)

type thread struct {
}

func (t *thread) ToConvert() (err error) {
	start := time.Now()

	cfg := mcfg.GetCfg()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxThreadTable := discuzPre + "forum_thread"
	dxPostTable := discuzPre + "forum_post"

	lastTid := cfg.GetInt("tables.xiuno.thread.last_tid")

	fields := "t.fid,t.tid,t.displayorder,t.authorid,t.subject,t.dateline,t.lastpost,t.views,t.replies,t.closed,p.useip,p.pid"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxThreadTable+" t").LeftJoin(dxPostTable+" p", "p.tid = t.tid").Where("p.first = ?", 1).Where("t.tid >= ?", lastTid).OrderBy("tid ASC").Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")
	if err != nil {
		mlog.Log.Debug("", "表 %s 数据查询失败, %s", xiunoTable, err.Error())
	}

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 无数据可以转换", xiunoTable)
		return nil
	}

	xiunoDB := database.GetXiunoDB()
	if _, err = xiunoDB.Exec("TRUNCATE " + xiunoTable); err != nil {
		return errors.New(fmt.Sprintf("清空数据表(%s)失败, %s", xiunoTable, err.Error()))
	}

	var count int64
	batch := cfg.GetInt("tables.xiuno.thread.batch")

	dataList := gdb.List{}
	for _, u := range r.ToList() {
		userip := common.Ip2long(gconv.String(u["useip"]))
		d := gdb.Map{
			"fid":         u["fid"],
			"tid":         u["tid"],
			"top":         u["displayorder"],
			"uid":         u["authorid"],
			"userip":      userip,
			"subject":     u["subject"],
			"create_date": u["dateline"],
			"last_date":   u["lastpost"],
			"views":       u["views"],
			"posts":       u["replies"],
			"closed":      u["closed"],
			"firstpid":    u["pid"],
		}

		// 批量插入数量
		if batch > 1 {
			dataList = append(dataList, d)
		} else {
			if res, err := xiunoDB.Insert(xiunoTable, d); err != nil {
				return errors.New(fmt.Sprintf("表 %s 数据插入失败, %s", xiunoTable, err.Error()))
			} else {
				c, _ := res.RowsAffected()
				count += c
			}
		}
	}

	if len(dataList) > 0 {
		if res, err := xiunoDB.BatchInsert(xiunoTable, dataList, batch); err != nil {
			return errors.New(fmt.Sprintf("表 %s 数据插入失败, %s", xiunoTable, err.Error()))
		} else {
			count, _ = res.RowsAffected()
		}
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))
	return
}

func NewThread() *thread {
	t := &thread{}
	return t
}
