package controllers

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/g/database/gdb"
	"time"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mlog"
)

type mypost struct {
}

func (t *mypost) ToConvert() (err error) {
	start := time.Now()

	cfg := mcfg.GetCfg()
	xiunoPre := database.GetPrefix("xiuno")

	xnThreadTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")

	fields := "uid,tid,pid"
	var r gdb.Result
	r, err = database.GetXiunoDB().Table(xnThreadTable).Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.mypost.name")
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
	if res, err := xiunoDB.Insert(xiunoTable, r); err != nil {
		return errors.New(fmt.Sprintf("表 %s 数据插入失败, %s", xiunoTable, err.Error()))
	} else {
		count, _ = res.RowsAffected()
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))
	return
}

func NewMypost() *mypost {
	t := &mypost{}
	return t
}
