package dx3ToXn4

import (
	"fmt"
	"github.com/skiy/golib"
	"log"
)

type forum struct {
	dxstr,
	xnstr dbstr
	count,
	total int
	dbname string
}

type forumFields struct {
	fid,
	name,
	rank,
	threads,
	todayposts,
	brief,
	announcement,
	seo_title,
	seo_keywords string
}

func (this *forum) update() {
	this.dbname = this.xnstr.DBPre + "forum"
	if !lib.AutoUpdate(this.xnstr.Auto, this.dbname) {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.dbname + " 失败: " + err.Error())
	}

	fmt.Printf("转换 %s 表成功，共(%d)条数据\r\n", this.dbname, count)
}

func (this *forum) toUpdate() (count int, err error) {
	dxpre := this.dxstr.DBPre

	dxtb1 := dxpre + "forum_forum"
	dxtb2 := dxpre + "forum_forumfield"

	fields := this.dxstr.FieldAddPrev("f1", "fid,name,rank,threads,todayposts")
	fields += "," + this.dxstr.FieldAddPrev("f2", "description,rules,seotitle,keywords")

	dxsql := fmt.Sprintf("SELECT %s FROM %s f1 LEFT JOIN %s f2 ON f2.fid = f1.fid WHERE f1.type = 'forum' AND status = 1 ORDER BY f1.fid ASC", fields, dxtb1, dxtb2)

	newFields := "fid,name,rank,threads,todayposts,brief,announcement,seo_title,seo_keywords"
	qmark := this.dxstr.FieldMakeQmark(newFields, "?")
	xnsql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", this.dbname, newFields, qmark)

	data, err := dxdb.Query(dxsql)
	if err != nil {
		log.Fatalln(dxsql, err.Error())
	}
	defer data.Close()

	xnClear := "TRUNCATE " + this.dbname
	_, err = xndb.Exec(xnClear)
	if err != nil {
		log.Fatalf(":::清空 %s 表失败: "+err.Error(), this.dbname)
	}
	fmt.Printf("清空 %s 表成功 \r\n", this.dbname)

	stmt, err := xndb.Prepare(xnsql)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %s 表\r\n", this.dbname)

	var field forumFields
	for data.Next() {
		err = data.Scan(
			&field.fid,
			&field.name,
			&field.rank,
			&field.threads,
			&field.todayposts,
			&field.brief,
			&field.announcement,
			&field.seo_title,
			&field.seo_keywords)

		_, err = stmt.Exec(
			&field.fid,
			&field.name,
			&field.rank,
			&field.threads,
			&field.todayposts,
			&field.brief,
			&field.announcement,
			&field.seo_title,
			&field.seo_keywords)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
			lib.UpdateProcess(fmt.Sprintf("正在升级第 %d 条数据", count), 0)
		}
	}

	if err = data.Err(); err != nil {
		log.Fatalf("user insert error: %s \r\n", err.Error())
	}

	return count, err
}
