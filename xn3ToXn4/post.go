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
	waitFix,
	count,
	total int
}

type postFields struct {
	tid, pid, uid, isfirst, create_date, userip, images, files, message, message_fmt string
}

func (this *post) update() {
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"post") {
		return
	}

	currentTime := time.Now()

	err := this.toUpdate(this.waitFix)
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "post 失败: " + err.Error())
	}

	fmt.Println("this.wait:", this.waitFix)
	this.toUpdate(this.waitFix)

	fmt.Printf("转换 %spost 表成功，共(%d)条数据\r\n", this.db3str.DBPre, this.count)

	fmt.Println("\r\n转换 post 表总耗时: ", time.Since(currentTime))
}

/**
unused
*/
func (this *post) toUpdateLess() (err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "tid,pid,uid,isfirst,create_date,userip,images,files,message,message_fmt"
	qmark := this.db3str.FieldMakeQmark(fields, "?")
	xn3 := fmt.Sprintf("SELECT %s FROM %spost ORDER BY pid ASC", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %spost (%s) VALUES (%s)", xn4pre, fields, qmark)

	xn3count := fmt.Sprintf("SELECT COUNT(*) AS count FROM %spost", xn3pre)
	rows := xiuno3db.QueryRow(xn3count)
	rows.Scan(&this.total)
	fmt.Printf("post total: %d \r\n", this.total)

	data, err := xiuno3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4Clear := "TRUNCATE `" + xn4pre + "post`"
	_, err = xiuno4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %spost 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %spost 表成功\r\n", xn4pre)

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
			this.count++
			lib.UpdateProcess(fmt.Sprintf("正在升级第 %d / %d 条 post", this.count, this.total), 0)
		}
	}

	if err = data.Err(); err != nil {
		log.Fatalln(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return err
}

func (this *post) toUpdate(fixFlag int) (err error) {

	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	query, err := xiuno3db.Query("select * from " + xn3pre + "post")
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

	qmark2 := this.db3str.FieldMakeQmark(fields, "?")
	xn4 := fmt.Sprintf("INSERT INTO %spost (%s) VALUES (%s)", xn4pre, fields, qmark2)
	//fmt.Println("Xiuno 3: " + xn3)
	//fmt.Println("Xiuno 5: " + xn5)

	if fixFlag == 1 {
		return this.fixPost(oldField, xn4, msgFmtExist)
	}

	xn3count := fmt.Sprintf("SELECT COUNT(*) AS count FROM %spost", xn3pre)
	rows := xiuno3db.QueryRow(xn3count)
	rows.Scan(&this.total)
	fmt.Printf("post total: %d \r\n", this.total)

	data, err := xiuno3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4Clear := "TRUNCATE `" + xn4pre + "post`"
	_, err = xiuno4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %spost 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %spost 表成功\r\n", xn4pre)

	fmt.Printf("正在升级 %spost 表\r\n", xn4pre)

	//dataArr := make([]postFields, ...)

	var dataArr []postFields
	var longDataArr, errLongDataArr [][]postFields

	var sqlStr string
	var sqlArr []string

	start := 0
	times := 0
	offset := 50
	maxTimes := 30
	errorCount := 0

	var field postFields
	for data.Next() {
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

			if field.message_fmt == "" {
				field.message_fmt = field.message
			}

			dataArr = append(dataArr, field)
			start++

			if start%offset == 0 {
				times++
				longDataArr = append(longDataArr, dataArr)
				dataArr = nil

				if times >= maxTimes {
					for _, v := range longDataArr {
						sqlArr = this.makeFileSql(qmark, v)
						sqlStr = xn5 + strings.Join(sqlArr, ",")
						_, err = xiuno4db.Exec(sqlStr)
						if err != nil {
							fmt.Printf("%d - v - 导入数据失败(%s) \r\n", start, err.Error())

							errLongDataArr = append(errLongDataArr, v)
							errorCount = len(errLongDataArr) * offset
						} else {
							this.count += len(v)
							lib.UpdateProcess(fmt.Sprintf("正在升级第 %d / %d 条 post，错误: %d", this.count, this.total, errorCount), 0)
						}
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

	//最后一次未满足 maxTimes 时, 剩下的数据
	if dataArr != nil {
		sqlArr = this.makeFileSql(qmark, dataArr)
		sqlStr = xn5 + strings.Join(sqlArr, ",")
		_, err = xiuno4db.Exec(sqlStr)

		if err != nil {
			fmt.Printf("dataArr - 导入数据失败(%s) \r\n", err.Error())

			errLongDataArr = append(errLongDataArr, dataArr)
			errorCount = len(errLongDataArr) * offset
		}
		this.count += len(dataArr)
		lib.UpdateProcess(fmt.Sprintf("正在升级第 %d / %d 条 post，错误: %d", this.count, this.total, errorCount), 0)
	}

	fmt.Println("errlongDataArr:", errLongDataArr)

	//处理错误部分的
	if errLongDataArr != nil {
		stmt, err := xiuno4db.Prepare(xn4)
		if err != nil {
			log.Fatalln("处理部分错误: " + err.Error())
		}

		start = 0
		errCount := 0
		for _, values := range errLongDataArr {
			for _, value := range values {
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
					errCount++
				} else {
					this.count++
					lib.UpdateProcess(fmt.Sprintf("正在升级第 %d / %d 条 post，错误: %d", this.count, this.total, errCount), 0)
				}
			}
		}
	}

	if err != nil {
		log.Fatalln("txErr: " + err.Error())
	}

	if this.total-this.count > 0 {
		//如果导入部分有失败的,则修复
		this.waitFix = 1
	}
	return err
}

func (this *post) fixPost(oldField, xn4 string, msgFmtExist bool) (err error) {
	sql := "SELECT " + oldField + " FROM %s WHERE pid NOT IN (SELECT pid FROM %s)"

	xn3dbName := this.db3str.DBName + "." + this.db3str.DBPre + "post"
	xn4dbName := this.db4str.DBName + "." + this.db4str.DBPre + "post"
	xn3sql := fmt.Sprintf(sql, xn3dbName, xn4dbName)

	data, err := xiuno3db.Query(xn3sql)
	if err != nil {
		fmt.Println("查询数据库失败", err.Error())
		return
	}
	data.Close()

	if data != nil {

		stmt, err := xiuno4db.Prepare(xn4)
		if err != nil {
			log.Fatalln("修复帖子部分错误: " + err.Error())
		}

		var field postFields
		for data.Next() {
			errCount := 0
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

			if field.message_fmt == "" {
				field.message_fmt = field.message
			}

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
				fmt.Printf("PID (%s) 导入数据失败(%s) \r\n", field.pid, err.Error())
				errCount++
			} else {
				this.count++
				lib.UpdateProcess(fmt.Sprintf("正在升级第 %d / %d 条 post，错误: %d", this.count, this.total, errCount), 0)
			}
		}
	}

	return
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
