package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skiy/xiuno-tools/lib"
	"log"
	"time"
)

type post struct {
	db3str,
	db4str dbstr
	fields postFields
}

type postFields struct {
	tid, pid, uid, isfirst, create_date, userip, images, files, message, message_fmt string
}

func (this *post) update() {
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"user_open_plat") {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "post 失败: " + err.Error())
	}

	fmt.Printf("转换 %spost 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *post) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "tid,pid,uid,isfirst,create_date,userip,images,files,message,message_fmt"
	qmark := this.db3str.FieldMakeQmark(fields)
	xn3 := fmt.Sprintf("SELECT %s FROM %spost", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %spost (%s) VALUES (%s)", xn4pre, fields, qmark)

	xn3db, err := this.db3str.Connect()
	data, err := xn3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4db, _ := this.db4str.Connect()
	xn4Clear := "TRUNCATE `" + xn4pre + "post`"
	_, err = xn4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %spost 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %spost 表成功\r\n", xn4pre)

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

	fmt.Printf("正在升级 %spost 表\r\n", xn4pre)
	for data.Next() {
		var field = this.fields
		err = data.Scan(
			&field.tid,
			&field.pid,
			&field.uid,
			&field.isfirst,
			&field.create_date,
			&field.userip,
			&field.images,
			&field.files,
			&field.message,
			&field.message_fmt)

		_, err = stmt.Exec(
			&field.tid,
			&field.pid,
			&field.uid,
			&field.isfirst,
			&field.create_date,
			&field.userip,
			&field.images,
			&field.files,
			&field.message,
			&field.message_fmt)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
		}

		xn4db.SetConnMaxLifetime(time.Second * 10)
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
