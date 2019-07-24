package extension

import (
	"errors"
	"fmt"
	"github.com/skiy/xiuno-tools/app/libraries/database"
	"github.com/skiy/xiuno-tools/app/libraries/mcfg"
	"github.com/skiy/xiuno-tools/app/libraries/mfile"
	"github.com/skiy/xiuno-tools/app/libraries/mlog"
	"github.com/skiy/xiuno-tools/app/libraries/mstr"
	"strings"
	"time"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/util/gconv"
)

// File 文件迁移
type File struct {
	DiscuzPath,
	XiunoPath string
}

// Parsing 解析
func (t *File) Parsing() (err error) {
	if !cfg.GetBool("extension.file.enable") {
		return
	}

	discuzPath := cfg.GetString("extension.file.discuzx_path")
	if discuzPath == "" {
		return errors.New("Discuz!X 站点路径 (discuzx_path) 未配置")
	}

	discuzPath = strings.TrimRight(discuzPath, "\\")
	t.DiscuzPath = strings.TrimRight(discuzPath, "/")
	if !gfile.IsDir(t.DiscuzPath) {
		return errors.New("Discuz!X 站点路径 (discuzx_path) 不是文件夹")
	}

	mlog.Log.Info("", "Discuz!X 的站点路径为: %s", t.DiscuzPath)

	xiunoPath := cfg.GetString("extension.file.xiuno_path")
	if xiunoPath == "" {
		// 工具当前目录
		xiunoPath = gfile.Pwd() + "/files"
		if err = gfile.Mkdir(xiunoPath); err != nil {
			err = fmt.Errorf("文件保存目录 (%s) 创建失败, %s", xiunoPath, err.Error())
			return
		}
		mlog.Log.Warning("", "XiunoBBS 站点路径 (xiuno_path) 未配置, 附件将移到至当前目录下的 files 目录下。转换成功后, 请将此目录下的 upload 复制到 XiunoBBS 根目录覆盖即可")
	}

	xiunoPath = strings.TrimRight(xiunoPath, "\\")
	t.XiunoPath = strings.TrimRight(xiunoPath, "/")
	if !gfile.IsDir(t.XiunoPath) {
		mlog.Log.Info("", "XiunoBBS 站点路径 (xiuno_path) 不是文件夹, 附件将移到至当前目录")
	}

	mlog.Log.Info("", "附件、头像、版块 ICON 将迁移至目录: %s", t.XiunoPath)

	if strings.EqualFold(t.XiunoPath, t.DiscuzPath) {
		return errors.New("Discuz!X 目录与附件迁移目录不能相同")
	}

	// 附件及图片迁移
	if cfg.GetBool("extension.file.attach") {
		if err := t.attachFiles(); err != nil {
			return err
		}
	}

	// 头像迁移
	if cfg.GetBool("extension.file.avatar") {
		if err := t.avatarImages(); err != nil {
			return err
		}
	}

	// 版块 ICON 迁移
	if cfg.GetBool("extension.file.icon") {
		if err := t.forumIcons(); err != nil {
			return err
		}
	}

	return nil
}

// NewFile File init
func NewFile() *File {
	t := &File{}
	return t
}

// attachFiles 迁移附件、图片文件
func (t *File) attachFiles() (err error) {
	xnAttachPath := t.XiunoPath + "/upload/attach"
	dxAttachPath := t.DiscuzPath + "/data/attachment/forum"

	if !gfile.IsDir(dxAttachPath) {
		return fmt.Errorf("Discuz!X 论坛附件目录 (%s) 不存在", dxAttachPath)
	}

	if err := gfile.Remove(xnAttachPath); err != nil {
		mlog.Log.Warning("", "迁移附件目录 (%s) 删除失败, %s", xnAttachPath, err.Error())
	}

	if err = mfile.CopyDir(dxAttachPath, xnAttachPath); err != nil {
		err = fmt.Errorf("\n迁移附件 (%s) \n至 (%s) 失败, \n原因: %s", dxAttachPath, xnAttachPath, err.Error())
		return
	}

	mlog.Log.Info("", "迁移附件、图片、文件操作成功")
	return nil
}

// avatarImages 迁移头像
func (t *File) avatarImages() (err error) {
	xnAvatarPath := t.XiunoPath + "/upload/avatar"
	dxAvatarPath := t.DiscuzPath + "/uc_server/data/avatar"

	if !gfile.IsDir(dxAvatarPath) {
		return fmt.Errorf("Discuz!X 论坛头像目录 (%s) 不存在", dxAvatarPath)
	}

	if err := gfile.Remove(xnAvatarPath); err != nil {
		mlog.Log.Warning("", "迁移头像目录 (%s) 删除失败, %s", xnAvatarPath, err.Error())
	}

	start := time.Now()

	cfg := mcfg.GetCfg()

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxMemberTable := discuzPre + "common_member"

	fields := "uid"
	var r gdb.Result

	w := g.Map{
		"avatarstatus": 1,
	}

	r, err = database.GetDiscuzDB().Table(dxMemberTable).Where(w).Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.user.name")
	if err != nil {
		return fmt.Errorf("表 %s 头像数据查询失败, %s", xiunoTable, err.Error())
	}

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 无头像数据可以迁移", xiunoTable)
		return nil
	}

	xiunoDB := database.GetXiunoDB()

	var count int64
	timestamp := time.Now().Unix()

	for _, u := range r.ToList() {
		uid := gconv.String(u["uid"])
		realUID := fmt.Sprintf("%09s", uid)

		// XiunoBBS avatar rule
		dir1 := mstr.SubStr(realUID, 0, 3)
		avatarImagesPath := fmt.Sprintf("%s/%s/", xnAvatarPath, dir1)
		xnAvatarFilePath := fmt.Sprintf("%s%s.png", avatarImagesPath, uid)

		if err = gfile.Mkdir(avatarImagesPath); err != nil {
			err = fmt.Errorf("头像保存目录 (%s) 创建失败, %s", avatarImagesPath, err.Error())
			return
		}

		// Discuz!X avatar rule
		dir2 := mstr.SubStr(realUID, 3, 2)
		dir3 := mstr.SubStr(realUID, 5, 2)
		dir4 := mstr.SubStr(realUID, -2, 0)
		dxAvatarImagePath := fmt.Sprintf("%s/%s/%s/%s/%s_avatar_big.jpg", dxAvatarPath, dir1, dir2, dir3, dir4)

		if !gfile.IsFile(dxAvatarImagePath) {
			err = fmt.Errorf("用户 UID (%s) 头像不存在: %s ", uid, dxAvatarImagePath)
			return
		}

		if err = mfile.CopyFile(dxAvatarImagePath, xnAvatarFilePath); err != nil {
			err = fmt.Errorf("\n迁移用户 (%s) 的头像 (%s) \n至 (%s) 失败, \n原因: %s", uid, dxAvatarImagePath, xnAvatarFilePath, err.Error())
			return
		}

		d := g.Map{
			"avatar": timestamp,
		}

		w := g.Map{
			"uid": uid,
		}

		res, err := xiunoDB.Table(xiunoTable).Data(d).Where(w).Update()
		if err != nil {
			return fmt.Errorf("表 %s 用户头像, , %s", xiunoTable, err.Error())
		}

		c, _ := res.RowsAffected()
		count += c
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 用户头像, 本次更新: %d 条数据, 耗时: %v", xiunoTable, count, time.Since(start)))
	return nil
}

// forumIcons
func (t *File) forumIcons() (err error) {
	xnIconPath := t.XiunoPath + "/upload/forum"
	dxIconPath := t.DiscuzPath + "/data/attachment/common"

	if !gfile.IsDir(dxIconPath) {
		return fmt.Errorf("Discuz!X 论坛版块 ICON 目录 (%s) 不存在", dxIconPath)
	}

	if err := gfile.Remove(xnIconPath); err != nil {
		mlog.Log.Warning("", "迁移版块 ICON 目录 (%s) 删除失败, %s", xnIconPath, err.Error())
	}

	if err := gfile.Mkdir(xnIconPath); err != nil {
		return fmt.Errorf("版块 ICON 目录 (%s) 创建失败, %s", xnIconPath, err.Error())
	}

	discuzPre, xiunoPre := database.GetPrefix("discuz"), database.GetPrefix("xiuno")

	dxForumField := discuzPre + "forum_forumfield"

	fields := "fid,icon"
	var r gdb.Result
	r, err = database.GetDiscuzDB().Table(dxForumField).Where("icon != ?", "").Fields(fields).Select()

	xiunoTable := xiunoPre + cfg.GetString("tables.xiuno.forum.name")
	if err != nil {
		return fmt.Errorf("表 %s 版块 ICON 数据查询失败, %s", xiunoTable, err.Error())
	}

	if len(r) == 0 {
		mlog.Log.Debug("", "表 %s 无版块 ICON 数据可以迁移", xiunoTable)
		return nil
	}

	xiunoDB := database.GetXiunoDB()

	timestamp := time.Now().Unix()

	for _, u := range r.ToList() {
		fid := gconv.Int(u["fid"])
		iconURL := gconv.String(u["icon"])

		xnIconPathFile := fmt.Sprintf("%s/%d.png", xnIconPath, fid)
		dxIconPathFile := fmt.Sprintf("%s/%s", dxIconPath, iconURL)

		if !gfile.IsFile(dxIconPathFile) {
			err = fmt.Errorf("版块 (%d) ICON不存在: %s ", fid, dxIconPathFile)
			return
		}

		if err = mfile.CopyFile(dxIconPathFile, xnIconPathFile); err != nil {
			err = fmt.Errorf("\n迁移版块 (%d) ICON (%s) \n至 (%s) 失败, \n原因: %s", fid, dxIconPath, xnIconPath, err.Error())
			return
		}

		d := g.Map{
			"icon": timestamp,
		}

		w := g.Map{
			"fid": fid,
		}

		if _, err := xiunoDB.Table(xiunoTable).Data(d).Where(w).Update(); err != nil {
			return fmt.Errorf("表 %s 版块 ICON, , %s", xiunoTable, err.Error())
		}
	}

	mlog.Log.Info("", fmt.Sprintf("表 %s 版块 ICON 更新成功", xiunoTable))
	return nil
}
