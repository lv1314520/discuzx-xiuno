package controllers

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/util/gconv"
	"time"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mlog"
)

type threadTop struct {
}

func (t *threadTop) ToConvert() (err error) {
	start := time.Now()

	cfg := mcfg.GetCfg()
	xiunoPre := database.GetPrefix("xiuno")

	xnThreadTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")

	fields := "tid,fid,top"
	var r gdb.Result
	r, err = database.GetXiunoDB().Table(xnThreadTable).Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.thread_top.name")
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
	dataList := gdb.List{}
	for _, u := range r.ToList() {
		if gconv.Int(u["top"]) == 0 {
			continue
		}

		dataList = append(dataList, u)
	}

	if res, err := xiunoDB.Insert(xiunoTable, dataList); err != nil {
		return errors.New(fmt.Sprintf("表 %s 数据插入失败, %s", xiunoTable, err.Error()))
	} else {
		count, _ = res.RowsAffected()
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))
	return
}

func NewThreadTop() *threadTop {
	t := &threadTop{}
	return t
}
