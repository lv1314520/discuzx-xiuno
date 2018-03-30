package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skiy/xiuno-tools/lib"
	"log"
	"strings"
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
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"post") {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "post 失败: " + err.Error())
	}

	fmt.Printf("转换 %spost 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

/**
unused
*/
func (this *post) toUpdateLess() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "tid,pid,uid,isfirst,create_date,userip,images,files,message,message_fmt"
	qmark := this.db3str.FieldMakeQmark(fields, "?")
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

	var field postFields
	for data.Next() {
		field = this.fields
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

		//xn4db.SetConnMaxLifetime(time.Second * 10)
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

func (this *post) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	xn3db, err := this.db3str.Connect()
	query, err := xn3db.Query("select * from " + xn3pre + "post")
	if err != nil {
		fmt.Println("查询数据库失败", err.Error())
		return
	}
	defer query.Close()

	oldField := "tid,pid,uid,isfirst,create_date,userip,images,files,message"
	fields := oldField + ",message_fmt"
	msgFmtExist := false
	cols, _ := query.Columns()
	for _, v := range cols {
		if v == "message_fmt" {
			oldField += ",message_fmt"
			msgFmtExist = true
			break
		}
	}

	xn3 := fmt.Sprintf("SELECT %s FROM %spost", oldField, xn3pre)
	xn5 := fmt.Sprintf("INSERT INTO %spost (%s) VALUES ", xn4pre, fields)
	qmark := this.db3str.FieldMakeValue(fields)

	//fmt.Println("Xiuno 3: " + xn3)
	//fmt.Println("Xiuno 5: " + xn5)

	data, err := xn3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4db, _ := this.db4str.Connect()

	xn4db.SetMaxIdleConns(0)
	xn4db.SetMaxOpenConns(100)

	xn4Clear := "TRUNCATE `" + xn4pre + "post`"
	_, err = xn4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %spost 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %spost 表成功\r\n", xn4pre)

	fmt.Printf("正在升级 %spost 表\r\n", xn4pre)

	//dataArr := make([]postFields, ...)

	var dataArr []postFields
	var longDataArr [][]postFields
	var sqlStr string
	var sqlArr []string

	var errLongDataArr [][]postFields

	start := 0
	times := 0
	offset := 100
	for data.Next() {
		var field postFields
		if msgFmtExist {
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
		} else {
			err = data.Scan(
				&field.tid,
				&field.pid,
				&field.uid,
				&field.isfirst,
				&field.create_date,
				&field.userip,
				&field.images,
				&field.files,
				&field.message)
		}

		if err != nil {
			fmt.Printf("获取数据失败(%s) \r\n", err.Error())
		} else {

			field.message = strings.Trim(field.message, "\\")
			field.message_fmt = strings.Trim(field.message_fmt, "\\")

			if field.message_fmt == "" {
				field.message_fmt = field.message
			}

			dataArr = append(dataArr, field)
			start++

			if start%offset == 0 {
				times++
				longDataArr = append(longDataArr, dataArr)
				dataArr = nil

				if times > 100 {
					for _, v := range longDataArr {
						sqlArr = this.makeFileSql(qmark, v)
						sqlStr = xn5 + strings.Join(sqlArr, ",")
						_, err = xn4db.Exec(sqlStr)
						if err != nil {
							fmt.Printf("%d.导入数据失败(%s) \r\n", start, err.Error())

							errLongDataArr = append(errLongDataArr, v)
							continue
						}
						count += len(v)

						lib.UpdateProcess(fmt.Sprintf("正在升级第 %d 条 post", count))
					}

					times = 0
					longDataArr = nil
				}

				start = 0
			}

		}
	}

	if err = data.Err(); err != nil {
		log.Fatalln("dataErr: " + err.Error())
	}

	if dataArr != nil {
		sqlArr = this.makeFileSql(qmark, dataArr)
		sqlStr = xn5 + strings.Join(sqlArr, ",")
		_, err = xn4db.Exec(sqlStr)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())

			errLongDataArr = append(errLongDataArr, dataArr)
		}
		count += len(dataArr)
	}

	//处理错误部分的
	if errLongDataArr != nil {
		qmark = this.db3str.FieldMakeQmark(fields, "?")
		xn4 := fmt.Sprintf("INSERT INTO %spost (%s) VALUES (%s)", xn4pre, fields, qmark)

		stmt, err := xn4db.Prepare(xn4)
		if err != nil {
			log.Fatalln("处理部分错误！" + err.Error())
		}

		start = 0
		for _, value := range errLongDataArr {
			for _, value := range value {
				start++
				fmt.Sprintf("插入错误序号: %d \r\n", start)

				_, err = stmt.Exec(
					&value.tid,
					&value.pid,
					&value.uid,
					&value.isfirst,
					&value.create_date,
					&value.userip,
					&value.images,
					&value.files,
					&value.message,
					&value.message_fmt)

				if err != nil {
					fmt.Printf("导入数据失败(%s) \r\n", err.Error())
				} else {
					count++
					lib.UpdateProcess(fmt.Sprintf("正在升级第 %d 条 post", count))
					xn4db.SetConnMaxLifetime(time.Second * 10)
				}
			}
		}
	}

	return count, err
}

func (this *post) makeFileSql(qmark string, dataArr []postFields) (dataStr []string) {
	for _, field := range dataArr {
		dataStr = append(dataStr, "("+fmt.Sprintf(qmark,
			field.tid,
			field.pid,
			field.uid,
			field.isfirst,
			field.create_date,
			field.userip,
			field.images,
			field.files,
			field.message,
			field.message_fmt)+")")
	}
	return
}
