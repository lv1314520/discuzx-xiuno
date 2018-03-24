package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type modlog struct {
	db3str,
	db4str dbstr
	fields modlogFields
}

type modlogFields struct {
	logid, uid, tid, pid, subject, comment, create_date, action string
}

func (this *modlog) update() {
	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "modlog 失败: " + err.Error())
	}

	fmt.Printf("转换 %smodlog 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *modlog) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "logid,uid,tid,pid,subject,comment,create_date,action"
	qmark := this.db3str.FieldMakeQmark(fields)
	xn3 := fmt.Sprintf("SELECT %s FROM %smodlog", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %smodlog (%s) VALUES (%s)", xn4pre, fields, qmark)

	xn3db, _ := this.db3str.Connect()
	data, err := xn3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4db, _ := this.db4str.Connect()
	xn4Clear := "TRUNCATE `" + xn4pre + "modlog`"
	_, err = xn4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %smodlog 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %smodlog 表成功\r\n", xn4pre)

	tx, err := xn4db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(xn4)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %smodlog 表\r\n", xn4pre)
	for data.Next() {
		var field = this.fields
		err = data.Scan(
			&field.logid,
			&field.uid,
			&field.tid,
			&field.pid,
			&field.subject,
			&field.comment,
			&field.create_date,
			&field.action)

		_, err = stmt.Exec(
			&field.logid,
			&field.uid,
			&field.tid,
			&field.pid,
			&field.subject,
			&field.comment,
			&field.create_date,
			&field.action)

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
