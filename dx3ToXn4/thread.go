package dx3ToXn4

import (
	"fmt"
	"github.com/skiy/golib"
	"log"
	"strconv"
)

type thread struct {
	dxstr,
	xnstr dbstr
	count,
	total int
	tbname string
}

type threadFields struct {
	fid,
	tid,
	top,
	uid,
	userip,
	subject,
	create_date,
	last_date,
	views,
	posts,
	closed,
	firstpid string
}

func (this *thread) update() {
	this.tbname = this.xnstr.DBPre + "thread"
	if !lib.AutoUpdate(this.xnstr.Auto, this.tbname) {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.tbname + " 失败: " + err.Error())
	}

	fmt.Printf("转换 %s 表成功，共(%d)条数据\r\n\r\n", this.tbname, count)
}

func (this *thread) toUpdate() (count int, err error) {
	dxpre := this.dxstr.DBPre
	xnpre := this.xnstr.DBPre

	dxtb1 := dxpre + "forum_thread"
	dxtb2 := dxpre + "forum_post"

	xntb2 := xnpre + "thread_top"
	xntb3 := xnpre + "mythread"

	fields := this.dxstr.FieldAddPrev("t", "fid,tid,displayorder,authorid,subject,dateline,lastpost,views,replies,closed")
	fields += "," + this.dxstr.FieldAddPrev("p", "useip,pid")

	dxsql := fmt.Sprintf("SELECT %s FROM %s t LEFT JOIN %s p ON p.tid = t.tid WHERE p.first = 1 ORDER BY t.tid ASC", fields, dxtb1, dxtb2)

	newFields := "fid,tid,top,uid,userip,subject,create_date,last_date,views,posts,closed,firstpid"
	qmark := this.dxstr.FieldMakeQmark(newFields, "?")
	xnsql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", this.tbname, newFields, qmark)

	xnsql2 := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", xntb2, "fid,tid,top", "?,?,?")
	xnsql3 := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", xntb3, "uid,tid", "?,?")

	data, err := dxdb.Query(dxsql)
	if err != nil {
		log.Fatalln(dxsql, err.Error())
	}
	defer data.Close()

	xnClear := "TRUNCATE " + this.tbname
	_, err = xndb.Exec(xnClear)
	if err != nil {
		log.Fatalf(":::清空 %s 表失败: "+err.Error(), this.tbname)
	}
	fmt.Printf("清空 %s 表成功 \r\n", this.tbname)

	xnClear2 := "TRUNCATE " + xntb2
	_, err = xndb.Exec(xnClear2)
	if err != nil {
		log.Fatalf(":::清空 %s 表失败: "+err.Error(), xntb2)
	}
	fmt.Printf("清空 %s 表成功 \r\n", xntb2)

	xnClear3 := "TRUNCATE " + xntb3
	_, err = xndb.Exec(xnClear3)
	if err != nil {
		log.Fatalf(":::清空 %s 表失败: "+err.Error(), xntb3)
	}
	fmt.Printf("清空 %s 表成功 \r\n", xntb3)

	stmt, err := xndb.Prepare(xnsql)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %s 表\r\n", this.tbname)

	var field threadFields
	for data.Next() {
		err = data.Scan(
			&field.fid,
			&field.tid,
			&field.top,
			&field.uid,
			&field.subject,
			&field.create_date,
			&field.last_date,
			&field.views,
			&field.posts,
			&field.closed,
			&field.userip,
			&field.firstpid)

		userip := lib.Ip2long(field.userip)

		_, err = stmt.Exec(
			&field.fid,
			&field.tid,
			&field.top,
			&field.uid,
			userip,
			&field.subject,
			&field.create_date,
			&field.last_date,
			&field.views,
			&field.posts,
			&field.closed,
			&field.firstpid)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
			lib.UpdateProcess(fmt.Sprintf("正在升级第 %d 条数据", count), 0)

			_, err = xndb.Exec(xnsql3, &field.uid, &field.tid)
			if err != nil {
				fmt.Printf("xnsql3 导入数据失败(%s) \r\n", err.Error())
			}

			top, _ := strconv.Atoi(field.top)
			if top > 0 {
				_, err = xndb.Exec(xnsql2, &field.fid, &field.tid, &field.top)
				if err != nil {
					fmt.Printf("xnsql2 导入数据失败(%s) \r\n", err.Error())
				}
			}
		}
	}

	if err = data.Err(); err != nil {
		log.Fatalf("user insert error: %s \r\n", err.Error())
	}

	return count, err
}
