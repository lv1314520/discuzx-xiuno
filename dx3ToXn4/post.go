package dx3ToXn4

import (
	"fmt"
	"github.com/skiy/bbcode"
	"github.com/skiy/golib"
	"log"
)

type post struct {
	dxstr,
	xnstr dbstr
	count,
	total int
	tbname string
}

type postFields struct {
	tid,
	pid,
	uid,
	isfirst,
	create_date,
	userip,
	message,
	message_fmt string
}

func (this *post) update() {
	this.tbname = this.xnstr.DBPre + "post"
	if !lib.AutoUpdate(this.xnstr.Auto, this.tbname) {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.tbname + " 失败: " + err.Error())
	}

	fmt.Printf("转换 %s 表成功，共(%d)条数据\r\n\r\n", this.tbname, count)
}

func (this *post) toUpdate() (count int, err error) {
	dxpre := this.dxstr.DBPre
	xnpre := this.xnstr.DBPre

	dxtb1 := dxpre + "forum_post"

	xntb2 := xnpre + "mypost"

	fields := "tid,pid,authorid,first,dateline,useip,message"

	dxsql := fmt.Sprintf("SELECT %s FROM %s ORDER BY pid ASC", fields, dxtb1)

	newFields := "tid,pid,uid,isfirst,create_date,userip,message,message_fmt"
	qmark := this.dxstr.FieldMakeQmark(newFields, "?")
	xnsql := fmt.Sprintf("INSERT INTO %s (%s, doctype) VALUES (%s, '0')", this.tbname, newFields, qmark)

	xnsql2 := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", xntb2, "uid,tid,pid", "?,?,?")

	data, err := dxdb.Query(dxsql)
	if err != nil {
		log.Fatalln(dxsql, err.Error())
	}
	defer data.Close()

	xnClear := "TRUNCATE " + this.tbname
	_, err = xndb.Exec(xnClear)
	if err != nil {
		log.Fatalf(":::清空 %s 表失败: "+err.Error(), this.tbname)
	}
	fmt.Printf("清空 %s 表成功 \r\n", this.tbname)

	xnClear2 := "TRUNCATE " + xntb2
	_, err = xndb.Exec(xnClear2)
	if err != nil {
		log.Fatalf(":::清空 %s 表失败: "+err.Error(), xntb2)
	}
	fmt.Printf("清空 %s 表成功 \r\n", xntb2)

	stmt, err := xndb.Prepare(xnsql)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %s 表\r\n", this.tbname)

	var field postFields
	var message_fmt string
	for data.Next() {
		err = data.Scan(
			&field.tid,
			&field.pid,
			&field.uid,
			&field.isfirst,
			&field.create_date,
			&field.userip,
			&field.message)

		userip := lib.Ip2long(field.userip)

		if field.message != "" {
			//message_fmt = lib.BBCodeToHtml(field.message) //未处理message中的附件的
			message_fmt = this.BBCodeToHtml(field.message) //处理message中的附件
		} else {
			message_fmt = ""
		}

		_, err = stmt.Exec(
			&field.tid,
			&field.pid,
			&field.uid,
			&field.isfirst,
			&field.create_date,
			userip,
			message_fmt,
			message_fmt)

		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
			lib.UpdateProcess(fmt.Sprintf("正在升级第 %d 条数据", count), 0)

			_, err = xndb.Exec(xnsql2, &field.uid, &field.tid, &field.pid)
			if err != nil {
				fmt.Printf("xnsql2 导入数据失败(%s) \r\n", err.Error())
			}
		}
	}

	if err = data.Err(); err != nil {
		log.Fatalf("user insert error: %s \r\n", err.Error())
	}

	return count, err
}

/**
bbcode 转 html
*/
func (this *post) BBCodeToHtml(msg string) string {
	compiler := bbcode.NewCompiler(true, true)

	//转 table
	compiler.SetTag("table", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "table"
		return out, true
	})

	//转 tr
	compiler.SetTag("tr", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "tr"

		return out, true
	})
	//转 td
	compiler.SetTag("td", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "td"

		return out, true
	})

	//ul
	compiler.SetTag("list", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "ul"

		return out, true
	})

	//text-align
	compiler.SetTag("align", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "div"
		value := node.GetOpeningTag().Value
		if value != "" {
			out.Attrs["style"] = "text-align: " + value
		}
		return out, true
	})

	//li -> 将 [*] 转为 li
	compiler.SetTag("*", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "li"

		return out, true
	})

	//font -> 清空 font
	compiler.SetTag("font", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = ""

		return out, true
	})

	//free -> 清空 free
	compiler.SetTag("free", nil)

	//hide -> 清空 hide
	compiler.SetTag("hide", nil)

	//处理message中的附件
	pre := this.xnstr.DBPre

	xntb1 := pre + "attach"
	selSql1 := "SELECT isimage,filename FROM %s WHERE aid = ?"
	xnsql1 := fmt.Sprintf(selSql1, xntb1)

	var isimage, filename string
	compiler.SetTag("attach", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {

		out := bbcode.NewHTMLTag("")
		out.Name = ""

		closeFlag := true

		value := node.GetOpeningTag().Value
		if value == "" {
			attachId := bbcode.CompileText(node)
			//fmt.Println("attachid:", attachId, "\r\n")

			if len(attachId) > 0 {
				err := xndb.QueryRow(xnsql1, attachId).Scan(&isimage, &filename)
				if err != nil {
					fmt.Printf("查询附件(%s)失败(%s) \r\n", attachId, err.Error())
				} else {
					if isimage == "1" {
						out.Name = "img"
						out.Attrs["src"] = "upload/attach/" + filename

						closeFlag = false
					} else {
						out.Name = "a"
						out.Attrs["href"] = "?attach-download-" + attachId + ".htm" //bbcode.ValidURL(filename)
						out.Attrs["target"] = "_blank"
						out.Value = "附件: "

						closeFlag = true
					}
				}
			}
		}
		return out, closeFlag
	})

	return compiler.Compile(msg)
}
