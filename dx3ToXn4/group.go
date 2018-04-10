package dx3ToXn4

import (
	"fmt"
	"github.com/skiy/golib"
	"log"
	"strconv"
)

type group struct {
	dxstr,
	xnstr dbstr
	count,
	total int
	dbname string
}

type groupFields struct {
	gid,
	name,
	creditsfrom,
	creditsto,
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
	allowviewip string
}

func (this *group) update() {
	this.dbname = this.xnstr.DBPre + "group"
	if !lib.AutoUpdate(this.xnstr.Auto, this.dbname) {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.dbname + " 失败: " + err.Error())
	}

	fmt.Printf("转换 %s 表成功，共(%d)条数据\r\n", this.dbname, count)
}

func (this *group) toUpdate() (count int, err error) {
	dxpre := this.dxstr.DBPre

	dxtb1 := dxpre + "common_usergroup"
	dxtb2 := dxpre + "common_usergroup_field"
	dxtb3 := dxpre + "common_admingroup"

	fields := this.dxstr.FieldAddPrev("u", "groupid,grouptitle,type,creditslower,creditshigher,allowvisit")
	fields += "," + this.dxstr.FieldAddPrev("f", "allowpost,allowreply,allowpostattach,allowgetattach")
	fields += "," + this.dxstr.FieldAddPrev("a", "allowstickthread,alloweditpost,allowdelpost,allowmovethread,allowbanvisituser,allowviewip")

	dxsql := fmt.Sprintf("SELECT %s FROM %s u LEFT JOIN %s f ON u.groupid = f.groupid LEFT JOIN %s a ON a.admingid = u.groupid ORDER BY u.groupid ASC", fields, dxtb1, dxtb2, dxtb3)

	newFields := "gid,name,creditsfrom,creditsto,allowread,allowthread,allowpost,allowattach,allowdown,allowtop,allowupdate,allowdelete,allowmove,allowbanuser,allowviewip"
	qmark := this.dxstr.FieldMakeQmark(newFields, "?")
	xnsql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", this.dbname, newFields, qmark)

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

	stmt, err := xndb.Prepare(xnsql)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %s 表\r\n", this.dbname)

	var field groupFields
	var gtype string
	for data.Next() {
		err = data.Scan(
			&field.gid,
			&field.name,
			&gtype,
			&field.creditsfrom,
			&field.creditsto,
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
			&field.allowviewip)

		if gtype == "member" {
			field.allowtop = "0"
			field.allowupdate = "0"
			field.allowdelete = "0"
			field.allowmove = "0"
			field.allowbanuser = "0"
			field.allowviewip = "0"
		}

		allowtop, _ := strconv.Atoi(field.allowtop)
		if allowtop > 0 {
			field.allowtop = "1"
		}

		_, err = stmt.Exec(
			&field.gid,
			&field.name,
			&field.creditsfrom,
			&field.creditsto,
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
			&field.allowviewip)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
			lib.UpdateProcess(fmt.Sprintf("正在升级第 %d 条数据", count), 0)
		}
	}

	if err = data.Err(); err != nil {
		log.Fatalf("user insert error: %s \r\n", err.Error())
	}

	return count, err
}
