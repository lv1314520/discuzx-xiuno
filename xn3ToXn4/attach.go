package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skiy/golib"
	"log"
)

type attach struct {
	db3str,
	db4str dbstr
	fields attachFields
}

type attachFields struct {
	aid, tid, pid, uid, filesize, width, height, filename, orgfilename, filetype, create_date, comment, downloads, credits, golds, rmbs string
}

func (this *attach) update() {
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"attach") {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "attach 失败: " + err.Error())
	}

	fmt.Printf("转换 %sattach 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *attach) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "aid,tid,pid,uid,filesize,width,height,filename,orgfilename,filetype,create_date,comment,downloads"
	qmark := this.db3str.FieldMakeQmark(fields, "?")
	xn3 := fmt.Sprintf("SELECT %s FROM %sattach", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %sattach (%s,credits,golds,rmbs) VALUES (%s, 0, 0, 0)", xn4pre, fields, qmark)

	data, err := xiuno3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4Clear := "TRUNCATE `" + xn4pre + "attach`"
	_, err = xiuno4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %sattach 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %sattach 表成功\r\n", xn4pre)

	tx, err := xiuno4db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(xn4)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %sattach 表\r\n", xn4pre)

	var field attachFields
	for data.Next() {
		err = data.Scan(
			&field.aid,
			&field.tid,
			&field.pid,
			&field.uid,
			&field.filesize,
			&field.width,
			&field.height,
			&field.filename,
			&field.orgfilename,
			&field.filetype,
			&field.create_date,
			&field.comment,
			&field.downloads)

		_, err = stmt.Exec(
			&field.aid,
			&field.tid,
			&field.pid,
			&field.uid,
			&field.filesize,
			&field.width,
			&field.height,
			&field.filename,
			&field.orgfilename,
			&field.filetype,
			&field.create_date,
			&field.comment,
			&field.downloads)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
		}
	}

	if err = data.Err(); err != nil {
		log.Fatalln(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return count, err
}
