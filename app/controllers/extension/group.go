package extension

import (
	"fmt"
	"strings"
	"time"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mlog"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/util/gconv"
)

/*
Group 用户组
*/
type Group struct {
}

/*
NewGroup Group init
*/
func NewGroup() *Group {
	t := &Group{}
	return t
}

/*
Parsing 解析
*/
func (t *Group) Parsing() (err error) {
	// 是否用户用户组变更
	if !cfg.GetBool("extension.group.enable") {
		return
	}

	// 使用官方组
	if cfg.GetBool("tables.xiuno.group.official") {
		return t.official()
	}

	return t.discuzGroup()
}

/**
官方组转换
*/
func (t *Group) official() (err error) {
	xiunoPre := database.GetPrefix("xiuno")
	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")

	var start time.Time
	var count int64
	var d g.Map

	d = g.Map{
		"gid": 101,
	}

	start = time.Now()
	r, err := database.GetXiunoDB().Table(xiunoTable).Data(d).Update()
	if err != nil {
		return fmt.Errorf("表 %s 重置用户组 gid 为 101 失败, %s", xiunoTable, err.Error())
	}
	count, _ = r.RowsAffected()
	mlog.Log.Info("", fmt.Sprintf("表 %s 重置用户组 gid 为 101 成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))

	// 不转换管理员 gid
	adminID := cfg.GetInt("extension.group.admin_id")
	if adminID <= 0 {
		return
	}

	d = g.Map{
		"gid": 1,
	}

	w := g.Map{
		"uid": adminID,
	}

	r, err = database.GetXiunoDB().Table(xiunoTable).Where(w).Data(d).Update()
	if err != nil {
		return fmt.Errorf("表 %s 重置 uid 为 %d 的用户组 gid 为 1 失败, %s", xiunoTable, adminID, err.Error())
	}

	count, _ = r.RowsAffected()
	mlog.Log.Info("", fmt.Sprintf("表 %s 重置 uid 为 %d 的用户组 gid 为 1 成功", xiunoTable, adminID))
	return
}

/**
Discuz 用户组数据修正
*/
func (t *Group) discuzGroup() (err error) {
	xiunoPre := database.GetPrefix("xiuno")
	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.group.name")

	fields := "gid,name"

	r, err := database.GetXiunoDB().Table(xiunoTable).Fields(fields).Select()
	if err != nil {
		return err
	}

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 无用户组可以转换", xiunoTable)
		return
	}

	guestMap := gmap.New()

	var tempGuessGid int
	var d, w g.Map
	var adminMap, powerArr g.ArrayStr

	for _, v := range r.ToList() {
		guestMap.Set(v["gid"], v["name"])

		// 内置游客组 ID
		if v["name"] == "游客" {
			tempGuessGid = gconv.Int(v["gid"])
		}

		// 可以删除用户的管理员组
		if v["name"] == "管理员" || v["name"] == "超级版主" {
			adminMap = append(adminMap, gconv.String(v["gid"]))
		}
	}

	// 游客组 gid
	guestGid := cfg.GetInt("extension.group.guest_gid")

	// 未配置游客组 ID, 则使用默认的
	if guestGid <= 0 {
		guestGid = tempGuessGid
	}

	// 若不存在此组则使用默认的 7
	if !guestMap.Contains(guestGid) {
		guestGid = 7
	}

	// 不存在此用户组
	if !guestMap.Contains(guestGid) {
		mlog.Log.Debug("", "表 %s 无用户组(%d)可以转换为游客组", xiunoTable, guestGid)
		return
	}

	d = g.Map{
		"gid": 0,
	}

	w = g.Map{
		"gid": guestGid,
	}

	if _, err := database.GetXiunoDB().Table(xiunoTable).Where(w).Data(d).Update(); err != nil {
		return fmt.Errorf("%s 原 %s 组(%d) 转换为游客组 gid 为 0 失败, %s", xiunoTable, guestMap.Get(guestGid), guestGid, err.Error())
	}
	mlog.Log.Info("", fmt.Sprintf("%s 原 %s 组(%d) 转换为游客组 gid 为 0 成功", xiunoTable, guestMap.Get(guestGid), guestGid))

	// 删除用户的权限
	deleteUserPower := cfg.GetString("extension.group.delete_user_power")
	if deleteUserPower != "" {
		powerArr = strings.Split(deleteUserPower, ",")
	}

	// 未设置权限组
	if len(powerArr) == 0 {
		powerArr = adminMap
	}

	// 原数据里也找不到管理员组, 则使用默认的 1,2
	if len(powerArr) == 0 {
		powerArr = g.ArrayStr{
			"1",
			"2",
		}
	}

	d = g.Map{
		"allowdeleteuser": 1,
	}

	_, err = database.GetXiunoDB().Table(xiunoTable).Where("gid IN (?)", powerArr).Data(d).Update()
	if err != nil {
		return fmt.Errorf("%s 用户组(%v)增加删除用户权限失败, %s", xiunoTable, powerArr, err.Error())
	}

	mlog.Log.Info("", fmt.Sprintf("%s 用户组(%v)增加删除用户权限成功", xiunoTable, powerArr))
	return
}
