package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skiy/xiuno-tools/lib"
	"log"
	"strconv"
)

type thread struct {
	db3str,
	db4str dbstr
	fields threadFields
}

type threadFields struct {
	fid, tid, top, uid, subject, create_date, last_date, views, posts, images, files, mods, closed, firstpid, lastuid, lastpid string
}

func (this *thread) update() {
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"thread") {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "thread 失败: " + err.Error())
	}

	fmt.Printf("转换 %sthread 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *thread) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "fid,tid,top,uid,subject,create_date,last_date,views,posts,images,files,mods,closed,firstpid,lastuid,lastpid"
	qmark := this.db3str.FieldMakeQmark(fields)
	xn3 := fmt.Sprintf("SELECT %s FROM %sthread", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %sthread (%s) VALUES (%s)", xn4pre, fields, qmark)
	xn4_1 := fmt.Sprintf("INSERT INTO %sthread_top SET fid=?, tid=?, top=?", xn4pre)
	xn4_2 := fmt.Sprintf("INSERT INTO %smythread SET uid=?, tid=?", xn4pre)

	xn3db, _ := this.db3str.Connect()
	data, err := xn3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4db, _ := this.db4str.Connect()
	xn4Clear := "TRUNCATE `" + xn4pre + "thread`"
	_, err = xn4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %sthread 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %sthread 表成功\r\n", xn4pre)

	xn4Clear1 := "TRUNCATE `" + xn4pre + "thread_top`"
	_, err = xn4db.Exec(xn4Clear1)
	if err != nil {
		log.Fatalf(":::清空 %sthread_top 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %sthread_top 表成功\r\n", xn4pre)

	xn4Clear2 := "TRUNCATE `" + xn4pre + "mythread`"
	_, err = xn4db.Exec(xn4Clear2)
	if err != nil {
		log.Fatalf(":::清空 %smythread 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %smythread 表成功\r\n", xn4pre)

	tx, err := xn4db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(xn4)
	if err != nil {
		log.Fatal(err)
	}

	stmt_1, err := tx.Prepare(xn4_1)
	defer stmt_1.Close()
	if err != nil {
		log.Fatal(err)
	}

	stmt_2, err := tx.Prepare(xn4_2)
	defer stmt_2.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("正在升级 %sthread 表\r\n", xn4pre)
	for data.Next() {
		var field = this.fields
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
			&field.images,
			&field.files,
			&field.mods,
			&field.closed,
			&field.firstpid,
			&field.lastuid,
			&field.lastpid)

		_, err = stmt.Exec(
			&field.fid,
			&field.tid,
			&field.top,
			&field.uid,
			&field.subject,
			&field.create_date,
			&field.last_date,
			&field.views,
			&field.posts,
			&field.images,
			&field.files,
			&field.mods,
			&field.closed,
			&field.firstpid,
			&field.lastuid,
			&field.lastpid)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			top, _ := strconv.Atoi(field.top)
			if top > 0 {
				_, err = stmt_1.Exec(&field.fid, &field.tid, &field.top)
				if err != nil {
					fmt.Printf("%sthread_top 导入数据失败(%s) \r\n", xn4pre, err.Error())
				}
			}

			_, err = stmt_2.Exec(&field.uid, &field.tid)
			if err != nil {
				fmt.Printf("%smythread 导入数据失败(%s) \r\n", xn4pre, err.Error())
			}
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
