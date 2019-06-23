package extension

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
	"time"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mlog"
)

type user struct {
}

func NewUser() *user {
	t := &user{}
	return t
}

func (t *user) Parsing() (err error) {
	// 修正 gid 为 101 的用户及用户组
	if cfg.GetBool("extension.user.normal_user") {
		if err := t.normalUser(); err != nil {
			return err
		}
	}

	// 修正用户主题和帖子统计
	if cfg.GetBool("extension.user.total") {
		err := t.threadPostStat()
		return err
	}
	return
}

/**
修正 gid 为 101 的用户及用户组
*/
func (t *user) normalUser() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoGroupTable := xiunoPre + cfg.GetString("tables.xiuno.group.name")
	xiunoUserTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")

	fields := "gid,name"
	var r gdb.Record
	r, err = database.GetXiunoDB().Table(xiunoGroupTable).Where("creditsfrom = ? AND creditsto > ?", 0, 0).Fields(fields).OrderBy("gid ASC").One()

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 无用户组可以转换", xiunoGroupTable)
		return
	}

	var d, w g.Map

	d = g.Map{
		"gid": 101,
	}

	w = g.Map{
		"gid": r["gid"],
	}

	if _, err := database.GetXiunoDB().Table(xiunoGroupTable).Where(w).Data(d).Update(); err != nil {
		return errors.New(fmt.Sprintf("%s 原 %v 组(%v) 转换为普通用户组 gid 为 101 失败, %s", xiunoGroupTable, r["name"], r["gid"], err.Error()))
	} else {
		mlog.Log.Info("", fmt.Sprintf("%s 原 %v 组(%v) 转换为普通用户组 gid 为 101 成功", xiunoGroupTable, r["name"], r["gid"]))
	}

	if res, err := database.GetXiunoDB().Table(xiunoUserTable).Where(w).Data(d).Update(); err != nil {
		return errors.New(fmt.Sprintf("%s 原 %v 组(%v)的用户转换为普通用户组 gid 为 101 失败, %s", xiunoGroupTable, r["name"], r["gid"], err.Error()))
	} else {
		count, _ := res.RowsAffected()
		mlog.Log.Info("", fmt.Sprintf("%s 原 %v 组(%v)的用户转换为普通用户组 gid 为 101 成功, 本次更新: %d 条数据", xiunoGroupTable, r["name"], r["gid"], count))
	}

	mlog.Log.Info("", fmt.Sprintf("修正 gid 为 101 的用户及用户组, 此次转换数据耗时: %v", time.Since(start)))
	return
}

/**
修正用户主题和帖子数量, 帖子包含主题和回复
*/
func (t *user) threadPostStat() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoUserTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")
	xiunoThreadTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")

	xiunoDB := database.GetXiunoDB()

	fields := "uid"
	var r gdb.Result
	r, err = xiunoDB.Table(xiunoUserTable).Fields(fields).Select()

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 无用户可以转换主题和帖子数量", xiunoUserTable)
		return
	}

	var count int64
	for _, u := range r.ToList() {
		w := g.Map{
			"uid": u["uid"],
		}
		posts, err := database.GetXiunoDB().Table(xiunoPostTable).Where(w).Fields("tid").Count()
		if err != nil {
			posts = 0
		}

		threads, err := database.GetXiunoDB().Table(xiunoThreadTable).Where(w).Fields("tid").Count()
		if err != nil {
			threads = 0
		}

		d := g.Map{
			"threads": threads,
			"posts":   posts,
		}

		if res, err := xiunoDB.Table(xiunoUserTable).Data(d).Where(w).Update(); err != nil {
			return errors.New(fmt.Sprintf("表 %s 用户帖子统计更新失败, %s", xiunoUserTable, err.Error()))
		} else {
			c, _ := res.RowsAffected()
			count += c
		}
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 用户帖子统计, 本次更新: %d 条数据, 耗时: %v", xiunoUserTable, count, time.Since(start)))
	return
}
