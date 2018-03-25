package xn3ToXn4

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skiy/xiuno-tools/lib"
	"log"
)

type friendlink struct {
	db3str,
	db4str dbstr
	fields friendlinkFields
}

type friendlinkFields struct {
	linkid,
	ftype,
	rank,
	create_date,
	name, url string
}

func (this *friendlink) update() {
	if !lib.AutoUpdate(this.db4str.Auto, this.db4str.DBPre+"friendlink") {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.db3str.DBPre + "friendlink 失败: " + err.Error())
	}

	fmt.Printf("转换 %sfriendlink 表成功，共(%d)条数据\r\n", this.db3str.DBPre, count)
}

func (this *friendlink) toUpdate() (count int, err error) {
	xn3pre := this.db3str.DBPre
	xn4pre := this.db4str.DBPre

	fields := "linkid,type,rank,create_date,name,url"
	qmark := this.db3str.FieldMakeQmark(fields, "?")
	xn3 := fmt.Sprintf("SELECT %s FROM %sfriendlink", fields, xn3pre)
	xn4 := fmt.Sprintf("INSERT INTO %sfriendlink (%s) VALUES (%s)", xn4pre, fields, qmark)

	createTable := `
CREATE TABLE IF NOT EXISTS %sfriendlink (
  linkid bigint(11) unsigned NOT NULL AUTO_INCREMENT,
  type smallint(11) NOT NULL DEFAULT '0',
  rank smallint(11) NOT NULL DEFAULT '0',
  create_date int(11) unsigned NOT NULL DEFAULT '0',
  name char(32) NOT NULL DEFAULT '',
  url char(64) NOT NULL DEFAULT '',
  PRIMARY KEY (linkid),
  KEY type (type)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8
`
	xn3db, _ := this.db3str.Connect()
	xn3CreateTable := fmt.Sprintf(createTable, xn3pre)
	_, err = xn3db.Exec(xn3CreateTable)
	if err != nil {
		log.Fatalln("Xiuno3: ", xn3CreateTable, err.Error())
	}

	data, err := xn3db.Query(xn3)
	if err != nil {
		log.Fatalln(xn3, err.Error())
	}
	defer data.Close()

	xn4db, _ := this.db4str.Connect()

	xn4CreateTable := fmt.Sprintf(createTable, xn4pre)
	_, err = xn3db.Exec(xn4CreateTable)
	if err != nil {
		log.Fatalln("Xiuno4: ", xn4CreateTable, err.Error())
	}

	xn4Clear := "TRUNCATE `" + xn4pre + "friendlink`"
	_, err = xn4db.Exec(xn4Clear)
	if err != nil {
		log.Fatalf(":::清空 %sfriendlink 表失败: "+err.Error(), xn4pre)
	}
	fmt.Printf("清空 %sfriendlink 表成功\r\n", xn4pre)

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

	fmt.Printf("正在升级 %sfriendlink 表\r\n", xn4pre)

	var field friendlinkFields
	for data.Next() {
		err = data.Scan(
			&field.linkid,
			&field.ftype,
			&field.rank,
			&field.create_date,
			&field.name,
			&field.url)

		_, err = stmt.Exec(
			&field.linkid,
			&field.ftype,
			&field.rank,
			&field.create_date,
			&field.name,
			&field.url)

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
