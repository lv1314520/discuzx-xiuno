package extension

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/util/gconv"
	"strings"
	"time"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mfile"
	"xiuno-tools/app/libraries/mlog"
	"xiuno-tools/app/libraries/mstr"
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

	t.DiscuzPath = strings.TrimRight(discuzPath, gfile.Separator) + gfile.Separator
	if !gfile.IsDir(t.DiscuzPath) {
		return errors.New("Discuz!X 站点路径 (discuzx_path) 不是文件夹")
	}

	mlog.Log.Info("", "Discuz!X 的站点路径为: %s", t.DiscuzPath)

	xiunoPath := cfg.GetString("extension.file.xiuno_path")
	if xiunoPath == "" {
		mlog.Log.Info("", "XiunoBBS 站点路径 (xiuno_path) 未配置, 附件将移到至当前目录")

		// 工具当前目录
		xiunoPath = gfile.Pwd() + gfile.Separator + "uploads"
		if err = gfile.Mkdir(xiunoPath); err != nil {
			err = fmt.Errorf("附件保存目录(%s)创建失败, %s", xiunoPath, err.Error())
			return
		}
	}

	t.XiunoPath = strings.TrimRight(xiunoPath, gfile.Separator) + gfile.Separator
	if !gfile.IsDir(t.XiunoPath) {
		mlog.Log.Info("", "XiunoBBS 站点路径 (xiuno_path) 不是文件夹, 附件将移到至当前目录")
	}

	mlog.Log.Info("", "附件将迁移至目录: %s", t.XiunoPath)

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

	return
}

// NewFile File init
func NewFile() *File {
	t := &File{}
	return t
}

// attachFiles 迁移附件、图片文件
func (t *File) attachFiles() (err error) {
	xnAttachPath := t.XiunoPath + "upload/attach/"
	dxAttachPath := t.DiscuzPath + "data/attachment/forum/"

	if !gfile.IsDir(dxAttachPath) {
		return fmt.Errorf("Discuz!X 论坛附件目录(%s)不存在", dxAttachPath)
	}

	if err := gfile.Remove(xnAttachPath); err != nil {
		mlog.Log.Warning("", "迁移附件目录(%s)删除失败, %s", xnAttachPath, err.Error())
	}

	if err = mfile.CopyDir(dxAttachPath, xnAttachPath); err != nil {
		err = fmt.Errorf("\n迁移附件 (%s) \n至 (%s) 失败, \n原因: %s", dxAttachPath, xnAttachPath, err.Error())
		return
	}

	mlog.Log.Debug("", "\n迁移附件 (%s) \n至 (%s) 成功", dxAttachPath, xnAttachPath)
	return nil
}

// avatarImages 迁移头像
func (t *File) avatarImages() (err error) {
	xnAvatarPath := t.XiunoPath + "upload/avatar"
	dxAvatarPath := t.DiscuzPath + "uc_server/data/avatar"

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
		mlog.Log.Debug("", "表 %s 头像数据查询失败, %s", xiunoTable, err.Error())
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
			err = fmt.Errorf("\n迁移用户(%s)的头像 (%s) \n至 (%s) 失败, \n原因: %s", uid, dxAvatarImagePath, xnAvatarFilePath, err.Error())
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
