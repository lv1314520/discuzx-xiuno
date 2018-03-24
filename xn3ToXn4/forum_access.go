package xn3ToXn4

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type forum_access struct {
	db3str,
	db4str dbstr
	fields forum_accessFields
}

type forum_accessFields struct {
	fid,gid,allowread,allowthread,allowpost,allowattach,allowdown string
}

func (this *forum_access) update() {
	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "forum_access 失败: " + err.Error())
	}

	fmt.Printf("转换 %sforum_access 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *forum_access) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "fid,gid,allowread,allowthread,allowpost,allowattach,allowdown"
	qmark := this.db3str.FieldMakeQmark(fields)
	xn3 := fmt.Sprintf("SELECT %s FROM %sforum_access", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %sforum_access (%s) VALUES (%s)", xn4pre, fields, qmark)

	xn3db, err := sql.Open("mysql", this.db3str.DSN)
	data, err := xn3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4db, _ := this.db4str.Connect()
	xn4Clear := "TRUNCATE `" + xn4pre + "forum_access`"
	_, err = xn4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %sforum_access 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %sforum_access 表成功\r\n", xn4pre)

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

	fmt.Printf("正在升级 %sforum_access 表\r\n", xn4pre)
	for data.Next() {
		var field = this.fields
		err = data.Scan(
			&field.fid,
&field.gid,
&field.allowread,
&field.allowthread,
&field.allowpost,
&field.allowattach,
&field.allowdown)

		_, err = stmt.Exec(
			&field.fid,
			&field.gid,
			&field.allowread,
			&field.allowthread,
			&field.allowpost,
			&field.allowattach,
			&field.allowdown)

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
