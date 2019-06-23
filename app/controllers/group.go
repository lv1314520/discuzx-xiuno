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

type group struct {
}

func (t *group) ToConvert() (err error) {
	cfg := mcfg.GetCfg()

	// 使用 XiunoBBS 官方用户组, 则不转换
	if cfg.GetBool("tables.xiuno.group.official") {
		return
	}

	start := time.Now()
	fmt.Println(start)

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxGroupTable := discuzPre + "common_usergroup"
	dxGroupField := discuzPre + "common_usergroup_field"
	dxAdminGroup := discuzPre + "common_admingroup"

	fields := "u.groupid,u.grouptitle,u.type,u.creditslower,u.creditshigher,u.allowvisit,f.allowpost,f.allowreply,f.allowpostattach,f.allowgetattach,a.allowstickthread,a.alloweditpost,a.allowdelpost,a.allowmovethread,a.allowbanvisituser,a.allowviewip"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxGroupTable+" u").LeftJoin(dxGroupField+" f", "f.groupid = u.groupid").LeftJoin(dxAdminGroup+" a", "a.admingid = u.groupid").Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.group.name")
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

	dataList := gdb.List{}
	for _, u := range r.ToList() {
		allowtop := gconv.Int(u["allowstickthread"])

		d := gdb.Map{
			"gid":          u["groupid"],
			"name":         u["grouptitle"],
			"creditsfrom":  u["creditshigher"],
			"creditsto":    u["creditslower"],
			"allowread":    u["allowvisit"],
			"allowthread":  u["allowpost"],
			"allowpost":    u["allowreply"],
			"allowattach":  u["allowpostattach"],
			"allowdown":    u["allowgetattach"],
			"allowtop":     allowtop,
			"allowupdate":  u["alloweditpost"],
			"allowdelete":  u["allowdelpost"],
			"allowmove":    u["allowmovethread"],
			"allowbanuser": u["allowbanvisituser"],
			"allowviewip":  u["allowviewip"],
		}

		// 普通会员
		if gconv.String(u["type"]) == "member" {
			d["allowtop"] = 0
			d["allowupdate"] = "0"
			d["allowdelete"] = "0"
			d["allowmove"] = "0"
			d["allowbanuser"] = "0"
			d["allowviewip"] = "0"
		}

		// 允许置顶,则值全为 1
		if gconv.Int(d["allowtop"]) > 0 {
			d["allowtop"] = 1
		}

		dataList = append(dataList, d)
	}

	if res, err := xiunoDB.BatchInsert(xiunoTable, dataList, 100); err != nil {
		return errors.New(fmt.Sprintf("表 %s 数据插入失败, %s", xiunoTable, err.Error()))
	} else {
		count, _ := res.RowsAffected()
		mlog.Log.Info("", fmt.Sprintf("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))
		return nil
	}

	return
}

func NewGroup() *group {
	t := &group{}
	return t
}