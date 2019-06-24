package extension

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gconv"
	"time"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mlog"
)

type threadPost struct {
}

func NewThreadPost() *threadPost {
	t := &threadPost{}
	return t
}

func (t *threadPost) Parsing() (err error) {
	// 是否修正主题的 lastpid 和 lastuid
	if cfg.GetBool("extension.thread_post.fix_last") {
		if err := t.fixThreadLast(); err != nil {
			return err
		}
	}

	// 修正帖子内附件统计数量
	if cfg.GetBool("extension.thread_post.post_attach_total") {
		if err := t.postAttachTotal(); err != nil {
			return err
		}
	}

	// 修正主题内附件统计数量
	if cfg.GetBool("extension.thread_post.thread_attach_total") {
		if err := t.threadAttachTotal(); err != nil {
			return err
		}
	}
	return
}

/**
是否修正主题的 lastpid 和 lastuid
*/
func (t *threadPost) fixThreadLast() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
	xiunoThreadable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")
	xiunoDB := database.GetXiunoDB()

	fields := "max(pid) as max_pid"
	r, err := xiunoDB.Table(xiunoPostTable).Fields(fields).GroupBy("tid").Select()
	if err != nil {
		return err
	}

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 无数据可以转换 lastpid 和 lastuid", xiunoPostTable)
		return
	}

	var pidArr g.ArrayStr
	for _, u := range r.ToList() {
		pidArr = append(pidArr, gconv.String(u["max_pid"]))
	}

	// 获取最后一条帖子的 tid,uid,pid
	fields_2 := "tid,pid,uid"
	res, err := xiunoDB.Table(xiunoPostTable).Where("pid in (?)", pidArr).Fields(fields_2).Select()
	if err != nil {
		return err
	}

	if len(res) == 0 {
		mlog.Log.Debug("", "表 %s 找不到 lastpid 和 lastuid", xiunoPostTable)
		return
	}

	var count int64
	for _, u := range res.ToList() {
		w := g.Map{
			"tid": u["tid"],
		}

		d := g.Map{
			"lastpid": u["pid"],
			"lastuid": u["uid"],
		}

		if res2, err := xiunoDB.Table(xiunoThreadable).Data(d).Where(w).Update(); err != nil {
			return errors.New(fmt.Sprintf("表 %s 更新帖子的 lastuid 和 lastuid 失败, %s", xiunoThreadable, err.Error()))
		} else {
			c, _ := res2.RowsAffected()
			count += c
		}
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 更新帖子的 lastuid 和 lastuid 成功, 本次更新: %d 条数据, 耗时: %v", xiunoThreadable, count, time.Since(start)))
	return
}

/**
修正主题内附件统计数量
*/
func (t *threadPost) threadAttachTotal() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
	xiunoThreadTable := xiunoPre + cfg.GetString("tables.xiuno.thread.name")
	xiunoDB := database.GetXiunoDB()

	var count int64
	r, err := xiunoDB.Table(xiunoThreadTable+" t").InnerJoin(xiunoPostTable+" p", "p.isfirst = 1 AND p.tid = t.tid").Data("t.files = p.files, t.images = p.images").Update()
	if err != nil {
		return errors.New(fmt.Sprintf("表 %s 更新主题的附件数(files)和图片数(images)失败, %s", xiunoThreadTable, err.Error()))
	}
	count, _ = r.RowsAffected()

	mlog.Log.Info("", fmt.Sprintf("表 %s 更新主题的附件数(files)和图片数(images)成功, 本次更新: %d 条数据, 耗时: %v", xiunoThreadTable, count, time.Since(start)))
	return
}

/**
修正帖子内附件统计数量
*/
func (t *threadPost) postAttachTotal() (err error) {
	start := time.Now()

	xiunoPre := database.GetPrefix("xiuno")
	xiunoAttachTable := xiunoPre + cfg.GetString("tables.xiuno.attach.name")
	xiunoPostTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
	xiunoDB := database.GetXiunoDB()

	fields := "count(*) as total, pid, isimage"
	r, err := xiunoDB.Table(xiunoAttachTable).Fields(fields).GroupBy("pid,isimage").Select()
	if err != nil {
		return err
	}

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 找不到附件数和图片数", xiunoAttachTable)
		return
	}

	var count int64
	var w, d g.Map
	for _, u := range r.ToList() {
		w = g.Map{
			"pid": u["pid"],
		}

		// 图片
		if u["isimage"] == 1 {
			d = g.Map{
				"images": u["total"],
			}
		} else { //非图片
			d = g.Map{
				"files": u["total"],
			}
		}

		if res, err := xiunoDB.Table(xiunoPostTable).Data(d).Where(w).Update(); err != nil {
			return errors.New(fmt.Sprintf("表 %s 更新帖子的附件数(files)和图片数(images)失败, %s", xiunoPostTable, err.Error()))
		} else {
			c, _ := res.RowsAffected()
			count += c
		}
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 更新帖子的附件数(files)和图片数(images)成功, 本次更新: %d 条数据, 耗时: %v", xiunoPostTable, count, time.Since(start)))
	return
}
