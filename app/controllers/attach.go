package controllers

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/util/gconv"
	"path"
	"strings"
	"time"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mlog"
)

type attach struct {
}

func (t *attach) ToConvert() (err error) {
	start := time.Now()

	cfg := mcfg.GetCfg()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxAttachmentTable := discuzPre + "forum_attachment"

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.attach.name")

	xiunoDB := database.GetXiunoDB()
	if _, err = xiunoDB.Exec("TRUNCATE " + xiunoTable); err != nil {
		return errors.New(fmt.Sprintf("清空数据表(%s)失败, %s", xiunoTable, err.Error()))
	}

	if err != nil {
		mlog.Log.Debug("", "表 %s 数据查询失败, %s", xiunoTable, err.Error())
	}

	var r gdb.Result
	var count int64
	var failureData []string

	fields := "aid,tid,pid,uid,filesize,width,attachment,filename,dateline,description,isimage"
	batch := cfg.GetInt("tables.xiuno.attach.batch")

	// forum_attachment_0 ~ forum_attachment_9
	for i := 0; i < 10; i++ {
		tbname := fmt.Sprintf("%s_%d", dxAttachmentTable, i)
		r, err = database.GetDiscuzDB().Table(tbname).Fields(fields).Select()

		if len(r) == 0 {
			mlog.Log.Debug("", "表 %s 无数据可以转换", xiunoTable)
			continue
		}

		dataList := gdb.List{}
		for _, u := range r.ToList() {
			filetype := t.FileExt(gconv.String(u["filename"]))

			isimage := gconv.Int(u["isimage"])
			if isimage != 1 {
				isimage = 0
			}

			d := gdb.Map{
				"aid":         u["aid"],
				"tid":         u["tid"],
				"pid":         u["pid"],
				"uid":         u["uid"],
				"filesize":    u["filesize"],
				"width":       u["width"],
				"filename":    u["attachment"],
				"orgfilename": u["filename"],
				"filetype":    filetype,
				"create_date": u["dateline"],
				"comment":     u["description"],
				"isimage":     isimage,
			}

			// 批量插入数量
			if batch > 1 {
				dataList = append(dataList, d)
			} else {
				if res, err := xiunoDB.Insert(xiunoTable, d); err != nil {
					//return errors.New(fmt.Sprintf("表 %s 数据插入失败, %s", xiunoTable, err.Error()))
					mlog.Log.Warning("", "表 %s 数据插入失败, %s", xiunoTable, err.Error())
					failureData = append(failureData, fmt.Sprintf("%s(aid:%v)", tbname, u["aid"]))
				} else {
					c, _ := res.RowsAffected()
					count += c
				}
			}
		}

		if len(dataList) > 0 {
			// 批量插入
			if res, err := xiunoDB.BatchInsert(xiunoTable, dataList, batch); err != nil {
				return errors.New(fmt.Sprintf("表 %s 数据插入失败, %s", xiunoTable, err.Error()))
			} else {
				c, _ := res.RowsAffected()
				count += c
			}
		}
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))

	if len(failureData) > 0 {
		mlog.Log.Warning("", "导入失败的数据: %v", failureData)
	}
	return nil
}

/**
获取文件后缀
*/
func (t *attach) FileExt(filename string) string {
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

func NewAttach() *attach {
	t := &attach{}
	return t
}
