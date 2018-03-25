package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skiy/xiuno-tools/lib"
	"log"
)

type group struct {
	db3str,
	db4str dbstr
	fields groupFields
}

type groupFields struct {
	gid,
	name,
	allowread,
	allowthread,
	allowpost,
	allowattach,
	allowdown,
	allowtop,
	allowupdate,
	allowdelete,
	allowmove,
	allowbanuser,
	allowdeleteuser,
	allowviewip string
}

func (this *group) update() {
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"group") {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "group 失败: " + err.Error())
	}

	fmt.Printf("转换 %sgroup 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *group) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "gid,name,allowread,allowthread,allowpost,allowattach,allowdown,allowtop,allowupdate,allowdelete,allowmove,allowbanuser,allowdeleteuser,allowviewip"
	qmark := this.db3str.FieldMakeQmark(fields, "?")
	xn3 := fmt.Sprintf("SELECT %s FROM %sgroup", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %sgroup (%s) VALUES (%s)", xn4pre, fields, qmark)

	xn3db, _ := this.db3str.Connect()
	data, err := xn3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4db, _ := this.db4str.Connect()
	xn4Clear := "TRUNCATE `" + xn4pre + "group`"
	_, err = xn4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %sgroup 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %sgroup 表成功\r\n", xn4pre)

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

	fmt.Printf("正在升级 %sgroup 表\r\n", xn4pre)

	for data.Next() {
		var field groupFields
		err = data.Scan(
			&field.gid,
			&field.name,
			&field.allowread,
			&field.allowthread,
			&field.allowpost,
			&field.allowattach,
			&field.allowdown,
			&field.allowtop,
			&field.allowupdate,
			&field.allowdelete,
			&field.allowmove,
			&field.allowbanuser,
			&field.allowdeleteuser,
			&field.allowviewip)

		_, err = stmt.Exec(
			&field.gid,
			&field.name,
			&field.allowread,
			&field.allowthread,
			&field.allowpost,
			&field.allowattach,
			&field.allowdown,
			&field.allowtop,
			&field.allowupdate,
			&field.allowdelete,
			&field.allowmove,
			&field.allowbanuser,
			&field.allowdeleteuser,
			&field.allowviewip)

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
