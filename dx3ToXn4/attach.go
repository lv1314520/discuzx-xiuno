package dx3ToXn4

import (
	"fmt"
	"github.com/skiy/golib"
	"log"
	"path"
	"strconv"
	"strings"
)

type attach struct {
	dxstr,
	xnstr dbstr
	count,
	total int
	tbname string
}

type attachFields struct {
	aid,
	tid,
	pid,
	uid,
	filesize,
	width,
	filename,
	orgfilename,
	filetype,
	create_date,
	comment,
	downloads,
	isimage string
}

func (this *attach) update() {
	this.tbname = this.xnstr.DBPre + "attach"
	if !lib.AutoUpdate(this.xnstr.Auto, this.tbname) {
		return
	}

	count, err := this.toUpdate()
	if err != nil {
		log.Fatalln("转换 " + this.tbname + " 失败: " + err.Error())
	}

	fmt.Printf("转换 %s 表成功，共(%d)条数据\r\n\r\n", this.tbname, count)
}

func (this *attach) toUpdate() (count int, err error) {
	dxpre := this.dxstr.DBPre

	dxtb1 := dxpre + "forum_attachment"

	fields := this.dxstr.FieldAddPrev("x", "aid,tid,pid,uid,filesize,width,attachment,filename,dateline,description,isimage")
	fields += "," + this.dxstr.FieldAddPrev("a", "downloads")

	//dxsql := fmt.Sprintf("SELECT %s FROM %s x LEFT JOIN %s a ON a.aid = x.aid  ORDER BY x.aid ASC", fields, dxtb1 + "_%s", dxtb1)

	newFields := "aid,tid,pid,uid,filesize,width,filename,orgfilename,filetype,create_date,comment,downloads,isimage"
	qmark := this.dxstr.FieldMakeQmark(newFields, "?")
	xnsql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", this.tbname, newFields, qmark)

	xnClear := "TRUNCATE " + this.tbname
	_, err = xndb.Exec(xnClear)
	if err != nil {
		log.Fatalf(":::清空 %s 表失败: "+err.Error(), this.tbname)
	}
	fmt.Printf("清空 %s 表成功 \r\n", this.tbname)

	stmt, err := xndb.Prepare(xnsql)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	fmt.Printf("正在升级 %s 表\r\n", this.tbname)

	var field attachFields
	var i int
	var filetype string

	for i = 0; i < 10; i++ {
		offset := strconv.Itoa(i)
		//dxsql = fmt.Sprintf(dxsql, offset)
		dxsql := fmt.Sprintf("SELECT %s FROM %s x LEFT JOIN %s a ON a.aid = x.aid  ORDER BY x.aid ASC", fields, dxtb1+"_"+offset, dxtb1)

		data, err := dxdb.Query(dxsql)
		if err != nil {
			log.Fatalln(dxsql, err.Error())
		}
		defer data.Close()

		for data.Next() {
			err = data.Scan(
				&field.aid,
				&field.tid,
				&field.pid,
				&field.uid,
				&field.filesize,
				&field.width,
				&field.filename,
				&field.orgfilename,
				&field.create_date,
				&field.comment,
				&field.isimage,
				&field.downloads)

			filetype = this.FileExt(field.orgfilename)

			if field.isimage != "1" {
				field.isimage = "0"
			}

			downloads := "0"
			if field.downloads != "" {
				downloads = field.downloads
			}

			_, err = stmt.Exec(
				&field.aid,
				&field.tid,
				&field.pid,
				&field.uid,
				&field.filesize,
				&field.width,
				&field.filename,
				&field.orgfilename,
				filetype,
				&field.create_date,
				&field.comment,
				downloads,
				&field.isimage)

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
	}

	return count, err
}

/**
获取文件后缀
*/
func (this *attach) FileExt(filename string) string {
	fileSuffix := path.Ext(filename)
	suffix := strings.Replace(fileSuffix, ".", "", -1)

	fileext := "other"

	if strings.EqualFold("png", suffix) || strings.EqualFold("jpg", suffix) ||
		strings.EqualFold("jpeg", suffix) || strings.EqualFold("bmp", suffix) {
		fileext = "image"
	} else if strings.EqualFold("rar", suffix) || strings.EqualFold("zip", suffix) {
		fileext = "zip"
	} else if strings.EqualFold("pdf", suffix) {
		fileext = "pdf"
	} else if strings.EqualFold("txt", suffix) {
		fileext = "text"
	} else if strings.EqualFold("xls", suffix) || strings.EqualFold("xlsx", suffix) ||
		strings.EqualFold("doc", suffix) || strings.EqualFold("docx", suffix) ||
		strings.EqualFold("ppt", suffix) || strings.EqualFold("pptx", suffix) {
		fileext = "office"
	}

	return fileext
}
