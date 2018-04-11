package dx3ToXn4

import (
	"fmt"
	"github.com/skiy/golib"
	"log"
)

type post struct {
	dxstr,
	xnstr dbstr
	count,
	total int
	dbname string
}

type postFields struct {
	tid,
	pid,
	uid,
	isfirst,
	create_date,
	userip,
	message,
	message_fmt string
}

func (this *post) update() {
	this.dbname = this.xnstr.DBPre + "post"
	if !lib.AutoUpdate(this.xnstr.Auto, this.dbname) {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.dbname + " 失败: " + err.Error())
	}

	fmt.Printf("转换 %s 表成功，共(%d)条数据\r\n", this.dbname, count)
}

func (this *post) toUpdate() (count int, err error) {
	dxpre := this.dxstr.DBPre
	xnpre := this.xnstr.DBPre

	dxtb1 := dxpre + "forum_post"

	xntb2 := xnpre + "mypost"

	fields := "tid,pid,authorid,first,dateline,useip,message"

	dxsql := fmt.Sprintf("SELECT %s FROM %s ORDER BY pid ASC", fields, dxtb1)

	newFields := "tid,pid,uid,isfirst,create_date,userip,message,message_fmt"
	qmark := this.dxstr.FieldMakeQmark(newFields, "?")
	xnsql := fmt.Sprintf("INSERT INTO %s (%s, doctype) VALUES (%s, '3')", this.dbname, newFields, qmark)

	xnsql2 := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", xntb2, "uid,tid,pid", "?,?,?")

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

	xnClear2 := "TRUNCATE " + xntb2
	_, err = xndb.Exec(xnClear2)
	if err != nil {
		log.Fatalf(":::清空 %s 表失败: "+err.Error(), xntb2)
	}
	fmt.Printf("清空 %s 表成功 \r\n", xntb2)

	stmt, err := xndb.Prepare(xnsql)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %s 表\r\n", this.dbname)

	var field postFields
	var message_fmt string
	for data.Next() {
		err = data.Scan(
			&field.tid,
			&field.pid,
			&field.uid,
			&field.isfirst,
			&field.create_date,
			&field.userip,
			&field.message)

		userip := lib.Ip2long(field.userip)

		if field.message != "" {
			message_fmt = lib.BBCodeToHtml(field.message)
		} else {
			message_fmt = ""
		}

		_, err = stmt.Exec(
			&field.tid,
			&field.pid,
			&field.uid,
			&field.isfirst,
			&field.create_date,
			userip,
			&field.message,
			message_fmt)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
			lib.UpdateProcess(fmt.Sprintf("正在升级第 %d 条数据", count), 0)

			_, err = xndb.Exec(xnsql2, &field.uid, &field.tid, &field.pid)
			if err != nil {
				fmt.Printf("xnsql2 导入数据失败(%s) \r\n", err.Error())
			}
		}
	}

	if err = data.Err(); err != nil {
		log.Fatalf("user insert error: %s \r\n", err.Error())
	}

	return count, err
}
