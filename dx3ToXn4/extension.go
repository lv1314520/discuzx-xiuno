package dx3ToXn4

import (
	"bufio"
	"fmt"
	"github.com/skiy/bbcode"
	"github.com/skiy/golib"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//group: ✔修正可删除用户的组 id,
//group: ✔将XiunoBBS 将creditsfrom为0，creditsto不为0的组ID改为101，并将 user 为此组的 gid 改为101
//post: ✔图片数及附件数从 attach 表中提取
//post: <err>message 中 [attach]1[/attach] 的内容提取并替换url - 会再将html给转换多一次，所以此功能提前
//thread: ✔图片数及附件数从 post 表中 isfirst = 1提取,
//thread: ✔修正最后发帖者及最后帖子
//user: ✔更新 threads 和 posts 统计

//Linux 平台下设置两论坛的根目录，移动附件、头像及版块图标
//user: ✔修正头像avatar
//forum: 修正版块icon
//forum: 修正版主uid

type extension struct {
	dxstr,
	xnstr dbstr
	count,
	total int
	tbname string
	dxpath,
	xnpath string
}

func (this *extension) update() {
	if !lib.AutoUpdate(this.xnstr.Auto, ", 修正最终数据") {
		return
	}

	//修正帖子图片 - 废弃
	//this.fixPostImages()

	//	//修正用户主题、帖子统计
	//	this.fixUserPostStat()
	//
	//	//修正用户组的删除用户权限
	//	this.fixGroup()
	//
	//	//修正gid为101的用户及用户组
	//	this.fixUserGroup()
	//
	//	//修正最后发帖者及最后帖子
	//	this.fixThreadLastPost()
	//
	//	//修正帖子的附件数和图片数
	//	this.fixPostAttach()
	//
	//	//修正主题的附件数和图片数
	//	this.fixThreadAttach()
	//
	//	//附件提示
	//	this.CopyAttachTip()

	buf := bufio.NewReader(os.Stdin)
	fmt.Println(`
----------------------------------
更新 版块icon、用户头像、版主
并且移动附件、版块icon、用户头像
是否更新其它扩展信息(Y/N): (默认为 N)
----------------------------------`)
	s := lib.Input(buf)
	if !strings.EqualFold(s, "Y") {
		return
	}

	//复制文件
	this.CopyFiles()
}

/**
内容中的附件与图片 (bug - 会再将html转换一次) 此功能废弃
*/
func (this *extension) fixPostImages() {
	pre := this.xnstr.DBPre

	this.tbname = pre + "post"
	xntb1 := pre + "attach"

	selSql := "SELECT pid,message,message_fmt FROM %s"
	xnsql := fmt.Sprintf(selSql, this.tbname)

	data, err := xndb.Query(xnsql)
	if err != nil {
		log.Fatalln(xnsql, err.Error())
	}
	defer data.Close()

	selSql1 := "SELECT isimage,filename FROM %s WHERE aid = ?"
	xnsql1 := fmt.Sprintf(selSql1, xntb1)

	upSql2 := "UPDATE %s SET message = ?, message_fmt = ? WHERE pid = ?"
	xnsql2 := fmt.Sprintf(upSql2, this.tbname)
	stmt, err := xndb.Prepare(xnsql2)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	var isimage, filename string

	compiler := bbcode.NewCompiler(true, true)
	compiler.SetTag("attach", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {

		out := bbcode.NewHTMLTag("")
		out.Name = ""

		closeFlag := true

		value := node.GetOpeningTag().Value
		if value == "" {
			attachId := bbcode.CompileText(node)
			//fmt.Println("attachid:", attachId, "\r\n")

			if len(attachId) > 0 {
				err = xndb.QueryRow(xnsql1, attachId).Scan(&isimage, &filename)
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

						closeFlag = true
					}
				}
			}
		}

		//fmt.Println(">>>>>>>>>>>>>>>>>\r\n", out ,"\r\n<<<<<<<<<<<<<<<<<<<<\r\n\r\n")

		return out, closeFlag
	})

	var pid, message, message_fmt string
	var count int
	for data.Next() {
		err = data.Scan(&pid, &message, &message_fmt)

		msg := compiler.Compile(message)
		_, err = stmt.Exec(&msg, &msg, &pid)
		if err != nil {
			fmt.Printf("导入数据失败(%s) \r\n", err.Error())
		} else {
			count++
			lib.UpdateProcess(fmt.Sprintf("正在更新 (post) 第 %d 条数据", count), 0)
		}
	}
	//
	//DefaultTagCompilers["url"] = func(node *BBCodeNode) (*HTMLTag, bool) {
	//	out := NewHTMLTag("")
	//	out.Name = "a"
	//	value := node.GetOpeningTag().Value
	//	if value == "" {
	//		text := CompileText(node)
	//		if len(text) > 0 {
	//			out.Attrs["href"] = ValidURL(text)
	//		}
	//	} else {
	//		out.Attrs["href"] = ValidURL(value)
	//	}
	//	return out, true
	//}
}

/**
修正用户主题与帖子统计
*/
func (this *extension) fixUserPostStat() {
	pre := this.xnstr.DBPre

	this.tbname = pre + "user"
	xntb1 := pre + "thread"
	xntb2 := pre + "post"

	upSql := `UPDATE %s u set u.threads = (SELECT count(*) FROM %s WHERE uid = u.uid), u.posts = (SELECT COUNT(*) FROM %s WHERE uid = u.uid)`
	xnsql := fmt.Sprintf(upSql, this.tbname, xntb1, xntb2)

	res, err := xndb.Exec(xnsql)
	if err != nil {
		errmsg := "更新用户主题、帖子统计失败: " + err.Error()
		fmt.Printf("error message: (%s) \r\n", errmsg)
	} else {
		rows, _ := res.RowsAffected()
		fmt.Printf("更新用户主题、帖子统计成功，共(%d)条数据\r\n\r\n", rows)
	}
}

func (this *extension) fixUserAvatar() {
	//var userAvatar bool
	//
	//buf := bufio.NewReader(os.Stdin)
	//fmt.Println("是否修正用户头像(Y/N): (默认为 N)")
	//s := lib.Input(buf)
	//if strings.EqualFold(s, "Y") {
	//	userAvatar = true
	//}
}

/**
更新用户组 (删除用户) 权限
*/
func (this *extension) fixGroup() {
	pre := this.xnstr.DBPre

	this.tbname = pre + "group"

	selSql := "SELECT gid,name FROM %s"
	xnsql1 := fmt.Sprintf(selSql, this.tbname)
	data, err := xndb.Query(xnsql1)
	if err != nil {
		log.Fatalln(xnsql1, err.Error())
	}
	defer data.Close()

	var gid, name string
	for data.Next() {
		err = data.Scan(&gid, &name)
		if err != nil {
			fmt.Printf("查询用户组失败(%s) \r\n", err.Error())
		} else {
			fmt.Printf("用户组ID：%s, 用户组名: %s \r\n", gid, name)
		}
	}

	buf := bufio.NewReader(os.Stdin)
	power := "1,2"
	var powerArr []string
	var powerList []int
	for {
		fmt.Println("为用户组添加(删除用户)权限，逗号隔开: (默认为 1,2)")

		s := lib.Input(buf)
		if s != "" {
			power = s
		}

		powerArr = strings.Split(power, ",")
		//拆分字符串并判断值是否为数字
		for _, p := range powerArr {
			val, err := strconv.Atoi(p)
			if err == nil {
				if val > 0 {
					powerList = append(powerList, val)
				}
			}
		}

		if len(powerList) > 0 {
			break
		}
	}

	if len(powerList) <= 0 {
		fmt.Printf("更新用户组(删除用户)权限失败 \r\n")
		return
	}

	powerStr := strings.Join(powerArr, ",")
	xnsql2 := fmt.Sprintf("UPDATE %s SET allowdeleteuser = 1 WHERE gid IN (%s)", this.tbname, powerStr)

	res, err := xndb.Exec(xnsql2)
	if err != nil {
		fmt.Printf("更新用户组(删除用户)权限失败: %s \r\n%s\r\n\r\n", err.Error())
	} else {
		rows, _ := res.RowsAffected()
		fmt.Printf("更新用户组(删除用户)权限成功，共(%d)条数据\r\n\r\n", rows)
	}
}

/**
更新最低级用户组ID为101,及对应用户的组为101
*/
func (this *extension) fixUserGroup() {
	pre := this.xnstr.DBPre

	this.tbname = pre + "group"

	selSql := "SELECT gid, name FROM %s WHERE creditsfrom = 0 AND creditsto > 0 ORDER BY gid ASC LIMIT 1"
	xnsql1 := fmt.Sprintf(selSql, this.tbname)
	var gid, name string
	err := xndb.QueryRow(xnsql1).Scan(&gid, &name)
	if err != nil {
		fmt.Printf("查询最初级用户组ID失败(%s) \r\n", err.Error())
		return
	}

	upSql := "UPDATE %s SET gid = 101 WHERE gid = %s"
	xnsql2 := fmt.Sprintf(upSql, this.tbname, gid)
	_, err = xndb.Exec(xnsql2)
	if err != nil {
		fmt.Printf("用户组 %s(%s) 修正为 101 失败(%s) \r\n", name, gid, err.Error())
		return
	} else {
		fmt.Printf("用户组 %s(%s) 修正为 101 成功\r\n\r\n", name, gid)
	}

	tbuser := pre + "user"
	upSql2 := "UPDATE %s SET gid = 101 WHERE gid = %s"
	xnsql3 := fmt.Sprintf(upSql2, tbuser, gid)
	res, err := xndb.Exec(xnsql3)
	if err != nil {
		fmt.Printf("用户组 %s(%s) 的用户修正为 101 失败(%s) \r\n", name, gid, err.Error())
		return
	} else {
		rows, _ := res.RowsAffected()
		fmt.Printf("用户组 %s(%s) 的用户修正为 101 成功，共(%d)条数据\r\n\r\n", name, gid, rows)
	}
}

/**
更新主题的 lastpid 和 lastuid
*/
func (this *extension) fixThreadLastPost() {
	pre := this.xnstr.DBPre

	this.tbname = pre + "thread"

	xntb1 := pre + "post"

	upsql := `
			UPDATE %s t
			INNER JOIN
			(
				SELECT tid, uid AS last_uid, pid AS last_pid
				FROM %s
				WHERE pid IN (SELECT max(pid) FROM %s GROUP BY tid)
			) p
				ON t.tid = p.tid
			SET
				t.lastuid = p.last_uid,
				t.lastpid = p.last_pid
			`
	xnsql := fmt.Sprintf(upsql, this.tbname, xntb1, xntb1)

	res, err := xndb.Exec(xnsql)
	if err != nil {
		fmt.Printf("更新主题的 lastpid 和 lastuid 失败", err.Error())
		return
	} else {
		rows, _ := res.RowsAffected()
		fmt.Printf("更新主题的 lastpid 和 lastuid 成功，共(%d)条数据\r\n\r\n", rows)
	}
}

/**
更新帖子的附件和图片数
*/
func (this *extension) fixPostAttach() {
	pre := this.xnstr.DBPre

	this.tbname = pre + "post"

	xntb1 := pre + "attach"

	upsql := `
			UPDATE %s p
			SET 
			images = (SELECT count(aid) FROM %s WHERE isimage = 1 AND pid = p.pid),
			files = (SELECT count(aid) FROM %s WHERE isimage != 1 AND pid = p.pid)
			
			`
	xnsql := fmt.Sprintf(upsql, this.tbname, xntb1, xntb1)

	res, err := xndb.Exec(xnsql)
	if err != nil {
		fmt.Printf("更新帖子附件(files)和图片数(images)失败", err.Error())
		return
	} else {
		rows, _ := res.RowsAffected()
		fmt.Printf("更新帖子附件(files)和图片数(images)成功，共(%d)条数据\r\n\r\n", rows)
	}
}

/**
更新主题的附件和图片数
*/
func (this *extension) fixThreadAttach() {
	pre := this.xnstr.DBPre

	this.tbname = pre + "thread"

	xntb1 := pre + "post"

	upsql := `
			UPDATE
				%s t
			INNER JOIN %s p ON
				p.isfirst = 1 AND p.tid = t.tid
			SET
				t.files = p.files,
				t.images = p.images
			`
	xnsql := fmt.Sprintf(upsql, this.tbname, xntb1)

	res, err := xndb.Exec(xnsql)
	if err != nil {
		fmt.Printf("更新主题附件(files)和图片数(images)失败", err.Error())
		return
	} else {
		rows, _ := res.RowsAffected()
		fmt.Printf("更新主题附件(files)和图片数(images)成功，共(%d)条数据\r\n\r\n", rows)
	}
}

/**
移动文件提示
*/
func (this *extension) CopyAttachTip() {
	fmt.Printf(`
请将 Discuz!X 的 data/attachment/forum/ 下的文件夹
复制到 XiunoBBS 的 upload/attach/ 中

`)
}

/**
复制文件
*/
func (this *extension) CopyFiles() {
	buf := bufio.NewReader(os.Stdin)
	for {
		if this.dxpath == "" {
			fmt.Println("\r\n配置 Discuz 根目录地址: ")

			s := lib.Input(buf)
			if s != "" {
				this.dxpath = strings.TrimSpace(s)
				this.dxpath = strings.TrimRight(this.dxpath, "/")
			}

			if this.dxpath == "" {
				continue
			}
		}

		if this.xnpath == "" {
			fmt.Println("\r\n配置 XiunoBBS 根目录地址: ")

			s := lib.Input(buf)
			if s != "" {
				this.xnpath = strings.TrimSpace(s)
				this.xnpath = strings.TrimRight(this.xnpath, "/")
			}

			if this.xnpath == "" {
				continue
			}
		}

		if this.dxpath == this.xnpath {
			fmt.Println("Discuz!X 和 XiunoBBS 目录地址不能相同")

			this.dxpath, this.xnpath = "", ""
			continue
		} else {
			fmt.Printf(`
Discuz!X 目录: %s
XiunoBBS 目录: %s

`, this.dxpath, this.xnpath)
			break
		}
	}

	//复制附件
	//this.copyAttachFiles()

	//复制头像
	//this.copyAvatarImages()

	//复制版块图标
	this.copyForumIcons()
}

/**
复制附件
*/
func (this *extension) copyAttachFiles() {
	if this.xnpath == "" || this.dxpath == "" {
		fmt.Printf(`
XiunoBBS 和 Discuz!X 目录不能为空
`)
		return
	}

	attachPath := this.xnpath + "/upload/attach"
	dxattachPath := this.dxpath + "/data/attachment/forum"

	buf := bufio.NewReader(os.Stdin)
	fmt.Printf(`
---------------------------------------
是否复制 Discuz 附件到 Xiuno
「注意」将清空以下文件夹，且不可恢复！ 
%s
请输入(Y/N): (默认为 N)
---------------------------------------
`, attachPath)
	s := lib.Input(buf)
	if !strings.EqualFold(s, "Y") {
		return
	}
	err := os.RemoveAll(attachPath)
	if err != nil {
		fmt.Printf(`
删除附件文件夹失败:
%s
errmsg: %s
`, attachPath, err.Error())
		return
	}

	err = lib.CopyDir(dxattachPath, attachPath)

	if err != nil {
		fmt.Printf(`
复制附件文件夹失败: 
%s 
-> 
%s
errmsg: %s
`, dxattachPath, attachPath, err.Error())

		return
	}

	fmt.Printf(`
复制附件文件夹成功: 
%s 
-> 
%s

`, dxattachPath, attachPath)

}

func (this *extension) copyAvatarImages() {
	if this.xnpath == "" || this.dxpath == "" {
		fmt.Printf(`
XiunoBBS 和 Discuz!X 目录不能为空
`)
		return
	}

	dxpre := this.dxstr.DBPre
	xnpre := this.xnstr.DBPre

	this.tbname = xnpre + "user"

	dxtb1 := dxpre + "common_member"

	selSql := "SELECT uid FROM %s WHERE avatarstatus = 1"

	dxsql := fmt.Sprintf(selSql, dxtb1)

	xnsql := fmt.Sprintf("UPDATE %s SET avatar = ? WHERE uid = ?", this.tbname)

	data, err := dxdb.Query(dxsql)
	if err != nil {
		log.Fatalln(dxsql, err.Error())
	}
	defer data.Close()

	stmt, err := xndb.Prepare(xnsql)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	avatarPath := this.xnpath + "/upload/avatar"
	dxavatarPath := this.dxpath + "/uc_server/data/avatar"

	var uid string
	var count int
	timestamp := time.Now().Unix()
	for data.Next() {
		err = data.Scan(&uid)

		if err != nil {
			fmt.Printf("获取用户头像失败(%s) \r\n", err.Error())
			continue
		}

		realUid := fmt.Sprintf("%09s", uid)

		//Xn avatar rule
		dir1 := lib.Substr(realUid, 0, 3)
		avatarImagePath := fmt.Sprintf("%s/%s/", avatarPath, dir1)
		avatarPathFile := avatarImagePath + uid + ".png"
		err = os.MkdirAll(avatarImagePath, os.ModePerm)
		if err != nil {
			fmt.Printf(`
创建用户(%s)头像文件夹失败: 
%s
errmsg: %s
`, uid, avatarImagePath, err.Error())

			continue
		}

		//Dx avatar rule
		dir2 := lib.Substr(realUid, 3, 2)
		dir3 := lib.Substr(realUid, 5, 2)
		dir4 := lib.Substr(realUid, -2, 0)
		dxAvatarImagePath := fmt.Sprintf("%s/%s/%s/%s/%s_avatar_big.jpg", dxavatarPath, dir1, dir2, dir3, dir4)

		_, err = os.Stat(dxAvatarImagePath)
		if err != nil {
			fmt.Printf(`
用户(%s)头像文件不存在: 
%s
errmsg: %s
`, uid, dxAvatarImagePath, err.Error())

			continue
		}

		err = lib.CopyFile(dxAvatarImagePath, avatarPathFile)

		if err != nil {
			fmt.Printf(`
复制用户(%s)头像失败: 
%s 
-> 
%s
errmsg: %s
`, uid, dxAvatarImagePath, avatarPathFile, err.Error())
		} else {

			_, err = xndb.Exec(xnsql, timestamp, uid)
			if err != nil {
				fmt.Printf("更新用户(%s)头像失败: %s", uid, err.Error())

				continue
			}

			count++
			lib.UpdateProcess(fmt.Sprintf("正在更新用户头像，第 %d 条数据", count), 0)
		}
	}

	/*
		$filename = "$uid.png";
		$dir = substr(sprintf("%09d", $uid), 0, 3).'/';
		$path = $conf['upload_path'].'avatar/'.$dir;
		$url = $conf['upload_url'].'avatar/'.$dir.$filename;

		$uid = abs(intval($uid));
		$uid = sprintf("%09d", $uid);
		$dir1 = substr($uid, 0, 3);
		$dir2 = substr($uid, 3, 2);
		$dir3 = substr($uid, 5, 2);
		$typeadd = $type == 'real' ? '_real' : '';
		return $dir1.'/'.$dir2.'/'.$dir3.'/'.substr($uid, -2).$typeadd."_avatar_$size.jpg";
	*/

	fmt.Printf("\r\n更新用户头像成功，共(%d)条数据\r\n\r\n", count)
}

/**
修正版块图标
*/
func (this *extension) copyForumIcons() {
	if this.xnpath == "" || this.dxpath == "" {
		fmt.Printf(`
XiunoBBS 和 Discuz!X 目录不能为空
`)
		return
	}

	dxpre := this.dxstr.DBPre
	xnpre := this.xnstr.DBPre

	this.tbname = xnpre + "forum"

	dxtb1 := dxpre + "forum_forumfield"

	selSql := "SELECT fid,icon FROM %s WHERE icon != ''"

	dxsql := fmt.Sprintf(selSql, dxtb1)

	xnsql := fmt.Sprintf("UPDATE %s SET icon = ? WHERE fid = ?", this.tbname)

	data, err := dxdb.Query(dxsql)
	if err != nil {
		log.Fatalln(dxsql, err.Error())
	}
	defer data.Close()

	stmt, err := xndb.Prepare(xnsql)
	if err != nil {
		log.Fatalf("stmt error: %s \r\n", err.Error())
	}
	defer stmt.Close()

	iconPath := this.xnpath + "/upload/forum"
	dxiconPath := this.dxpath + "/data/attachment/common"

	var fid, iconUrl string
	var count int
	timestamp := time.Now().Unix()
	for data.Next() {
		err = data.Scan(&fid, &iconUrl)

		if err != nil {
			fmt.Printf("获取版块(%s)icon失败(%s) \r\n", fid, err.Error())
			continue
		}

		iconPathFile := fmt.Sprintf("%s/%s.png", iconPath, fid)
		dxiconPathFile := fmt.Sprintf("%s/%s", dxiconPath, iconUrl)

		_, err = os.Stat(dxiconPathFile)
		if err != nil {
			fmt.Printf(`
版块(%s)icon文件不存在: 
%s
errmsg: %s
`, fid, dxiconPathFile, err.Error())

			continue
		}

		err = lib.CopyFile(dxiconPathFile, iconPathFile)

		if err != nil {
			fmt.Printf(`
版块(%s)icon文件失败: 
%s 
-> 
%s
errmsg: %s
`, fid, dxiconPathFile, iconPathFile, err.Error())
		} else {

			_, err = xndb.Exec(xnsql, timestamp, fid)
			if err != nil {
				fmt.Printf("更新版块(%s)icon失败: %s", fid, err.Error())

				continue
			}

			count++
			lib.UpdateProcess(fmt.Sprintf("正在更新版块icon，第 %d 条数据", count), 0)
		}
	}

	fmt.Printf("\r\n更新版块icon成功，共(%d)条数据\r\n\r\n", count)
}
