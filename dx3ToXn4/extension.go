package dx3ToXn4

import (
	"bufio"
	"fmt"
	"github.com/skiy/golib"
	"log"
	"os"
	"strconv"
	"strings"
)

//group: ✔修正可删除用户的组 id,
//group: ✔将XiunoBBS 将creditsfrom为0，creditsto不为0的组ID改为101，并将 user 为此组的 gid 改为101
//post: ✔图片数及附件数从 attach 表中提取
//post: message 中 [attach]1[/attach] 的内容提取并替换url
//thread: ✔图片数及附件数从 post 表中 isfirst = 1提取,
//thread: ✔修正最后发帖者及最后帖子
//user: ✔更新 threads 和 posts 统计
//forum: 修正版主 UID

//Linux 平台下设置两论坛的根目录，移动附件、头像及版块图标
//user: 修正头像avatarstatus
//forum: 修正版块icon

type extension struct {
	dxstr,
	xnstr dbstr
	count,
	total int
	tbname string
}

func (this *extension) update() {
	if !lib.AutoUpdate(this.xnstr.Auto, ", 修正最终数据") {
		return
	}

	//修正用户主题、帖子统计
	this.fixUserPostStat()

	//修正用户组的删除用户权限
	this.fixGroup()

	//修正gid为101的用户及用户组
	this.fixUserGroup()

	//修正最后发帖者及最后帖子
	this.fixThreadLastPost()

	//修正帖子的附件数和图片数
	this.fixPostAttach()

	//修正主题的附件数和图片数
	this.fixThreadAttach()

	//附件提示
	this.CopyAttachTip()

	buf := bufio.NewReader(os.Stdin)
	fmt.Println(`
是否更新其它扩展信息(Y/N): (默认为 N)
更新 版块icon、用户头像、版主
并且移动附件、版块icon、用户头像

`)
	s := lib.Input(buf)
	if !strings.EqualFold(s, "Y") {
		return
	}

	//复制文件
	//this.CopyFiles()
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

	xntb2 := pre + "post"

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
	xnsql := fmt.Sprintf(upsql, this.tbname, xntb2, xntb2)

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

	xntb2 := pre + "attach"

	upsql := `
			UPDATE %s p
			SET 
			images = (SELECT count(aid) FROM %s WHERE isimage = 1 AND pid = p.pid),
			files = (SELECT count(aid) FROM %s WHERE isimage != 1 AND pid = p.pid)
			
			`
	xnsql := fmt.Sprintf(upsql, this.tbname, xntb2, xntb2)

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

	xntb2 := pre + "post"

	upsql := `
			UPDATE
				%s t
			INNER JOIN %s p ON
				p.isfirst = 1 AND p.tid = t.tid
			SET
				t.files = p.files,
				t.images = p.images
			`
	xnsql := fmt.Sprintf(upsql, this.tbname, xntb2)

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
	var dxpath, xnpath string
	for {
		if dxpath == "" {
			fmt.Println("配置 Discuz 根目录地址: ")

			s := lib.Input(buf)
			if s != "" {
				dxpath = s
			}

			if dxpath == "" {
				continue
			}
		}

		if xnpath == "" {
			fmt.Println("配置 XiunoBBS 根目录地址: ")

			s := lib.Input(buf)
			if s != "" {
				xnpath = s
			}

			if xnpath == "" {
				continue
			}
		}

		if dxpath == xnpath {
			fmt.Println("Discuz 和 XiunoBBS 目录地址不能相同")

			dxpath, xnpath = "", ""
			continue
		} else {
			fmt.Printf("Discuz目录: %s\r\nXiunoBBS目录: %s\r\n\r\n", dxpath, xnpath)
			break
		}
	}

	err := lib.CopyDir(dxpath, xnpath)

	if err != nil {
		fmt.Println(err)
		return
	}
}

/**
修正版块图标
*/
func (this *extension) fixForumIcon() {
	//data/attachment/common/a5/common_{$fid}_icon.png - upload/forum/{$fid}.png
}
