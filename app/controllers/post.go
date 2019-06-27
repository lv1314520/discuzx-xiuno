package controllers

import (
	"fmt"
	"time"
	"xiuno-tools/app/libraries/common"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mlog"

	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/util/gconv"
	"github.com/skiy/bbcode"
)

type post struct {
}

func (t *post) ToConvert() (err error) {
	start := time.Now()

	cfg := mcfg.GetCfg()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxPostTable := discuzPre + "forum_post"

	lastPid := cfg.GetInt("tables.xiuno.post.last_pid")

	fields := "tid,pid,authorid,first,dateline,useip,message"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxPostTable+" t").Where("pid >= ?", lastPid).OrderBy("pid ASC").Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.post.name")
	if err != nil {
		mlog.Log.Debug("", "表 %s 数据查询失败, %s", xiunoTable, err.Error())
	}

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 无数据可以转换", xiunoTable)
		return nil
	}

	xiunoDB := database.GetXiunoDB()
	if _, err = xiunoDB.Exec("TRUNCATE " + xiunoTable); err != nil {
		return fmt.Errorf("清空数据表(%s)失败, %s", xiunoTable, err.Error())
	}

	var count int64
	batch := cfg.GetInt("tables.xiuno.post.batch")

	dataList := gdb.List{}
	for _, u := range r.ToList() {
		userip := common.IP2Long(gconv.String(u["useip"]))
		message_fmt := gconv.String(u["message"])

		if message_fmt != "" {
			message_fmt = t.BBCodeToHtml(message_fmt) //处理message中的附件
		}

		d := gdb.Map{
			"tid":         u["tid"],
			"pid":         u["pid"],
			"uid":         u["authorid"],
			"isfirst":     u["first"],
			"create_date": u["dateline"],
			"userip":      userip,
			"message":     message_fmt,
			"message_fmt": message_fmt,
		}

		// 批量插入数量
		if batch > 1 {
			dataList = append(dataList, d)
		} else {
			if res, err := xiunoDB.Insert(xiunoTable, d); err != nil {
				return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
			} else {
				c, _ := res.RowsAffected()
				count += c
			}
		}
	}

	if len(dataList) > 0 {
		if res, err := xiunoDB.BatchInsert(xiunoTable, dataList, batch); err != nil {
			return fmt.Errorf("表 %s 数据插入失败, %s", xiunoTable, err.Error())
		} else {
			count, _ = res.RowsAffected()
		}
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 数据导入成功, 本次导入: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))
	return
}

func NewPost() *post {
	t := &post{}
	return t
}

/**
bbcode 转 html
*/
func (t *post) BBCodeToHtml(msg string) string {
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

	//backcolor=yellow
	compiler.SetTag("backcolor", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "span"
		value := node.GetOpeningTag().Value
		if value != "" {
			out.Attrs["style"] = "background-color: " + value
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
	compiler.SetTag("free", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = ""

		return out, true
	})

	//hide -> 清空 hide
	compiler.SetTag("hide", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = ""

		return out, true
	})

	//qq -> 更新 qq 标签
	compiler.SetTag("qq", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = ""

		value := node.GetOpeningTag().Value
		if value == "" {
			qq := bbcode.CompileText(node)

			if len(qq) > 0 {
				out.Name = "a"
				out.Attrs["href"] = "http://wpa.qq.com/msgrd?v=3&uin=" + qq + "&site=Xiuno&from=Xiuno&menu=yes"
				out.Attrs["target"] = "_blank"
			}
		}
		return out, true
	})

	//处理message中的附件
	xiunoTable := database.GetPrefix("xiuno") + mcfg.GetCfg().GetString("tables.xiuno.attach.name")

	compiler.SetTag("attach", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {

		out := bbcode.NewHTMLTag("")
		out.Name = ""

		closeFlag := true

		value := node.GetOpeningTag().Value
		if value == "" {
			attachId := bbcode.CompileText(node)

			if len(attachId) > 0 {
				r, err := database.GetXiunoDB().Table(xiunoTable).Where("aid = ?", attachId).Fields("isimage,filename").One()

				if err != nil {
					mlog.Log.Warning("", "查询附件(aid: %s)失败, %s", attachId, err.Error())
				} else if r != nil {

					isimage := r["isimage"].Int()
					if isimage == 1 {
						out.Name = "img"
						out.Attrs["src"] = "upload/attach/" + r["filename"].String()

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
