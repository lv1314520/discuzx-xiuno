package controllers

import (
	"fmt"
	"github.com/skiy/xiuno-tools/app/libraries/database"
	"github.com/skiy/xiuno-tools/app/libraries/mcfg"
	"github.com/skiy/xiuno-tools/app/libraries/mlog"
	"time"

	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/util/gconv"
)

// ThreadTop 置顶
type ThreadTop struct {
}

// ToConvert 转换
func (t *ThreadTop) ToConvert() (err error) {
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
		return fmt.Errorf("清空数据表(%s)失败, %s", xiunoTable, err.Error())
	}

	var count int64
	dataList := gdb.List{}
	for _, u := range r.ToList() {
		top := gconv.Int(u["top"])
		if top != 1 && top != 2 && top != 3 {
			continue
		}
		dataList = append(dataList, u)
	}

	if len(dataList) == 0 {
		mlog.Log.Debug("", "表 %s 无数据可以转换", xiunoTable)
		return nil
	}

	res, err := xiunoDB.Insert(xiunoTable, dataList)
	if err != nil {
		return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
	}
	count, _ = res.RowsAffected()

	mlog.Log.Info("", fmt.Sprintf("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))
	return
}

// NewThreadTop ThreadTop init
func NewThreadTop() *ThreadTop {
	t := &ThreadTop{}
	return t
}
