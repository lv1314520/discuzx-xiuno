package controllers

import (
	"fmt"
	"github.com/skiy/xiuno-tools/app/libraries/database"
	"github.com/skiy/xiuno-tools/app/libraries/mcfg"
	"github.com/skiy/xiuno-tools/app/libraries/mlog"
	"time"

	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
)

// Forum 版块
type Forum struct {
}

// ToConvert 转换
func (t *Forum) ToConvert() (err error) {
	start := time.Now()

	cfg := mcfg.GetCfg()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxForumTable := discuzPre + "forum_forum"
	dxForumField := discuzPre + "forum_forumfield"

	fields := "f.fid,f.name,f.rank,f.threads,f.todayposts,e.description,e.rules,e.seotitle,e.keywords"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxForumTable+" f").LeftJoin(dxForumField+" e", "e.fid = f.fid").Where("f.type = ? AND f.status = ?", "forum", 1).Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.forum.name")
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
		d := gdb.Map{
			"fid":          u["fid"],
			"name":         u["name"],
			"rank":         u["rank"],
			"threads":      u["threads"],
			"todayposts":   u["todayposts"],
			"brief":        u["description"],
			"announcement": u["rules"],
		}

		seoTitle := gstr.SubStr(gconv.String(u["seotitle"]), 0, 64)
		seoKeywords := gstr.SubStr(gconv.String(u["keywords"]), 0, 64)

		d["seo_title"] = seoTitle
		d["seo_keywords"] = seoKeywords

		dataList = append(dataList, d)
	}

	if len(dataList) > 0 {
		res, err := xiunoDB.BatchInsert(xiunoTable, dataList, 100)
		if err != nil {
			return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
		}
		count, _ = res.RowsAffected()
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))
	return
}

// NewForum Forum init
func NewForum() *Forum {
	t := &Forum{}
	return t
}
