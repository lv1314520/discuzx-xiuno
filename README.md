### Xiuno Tools
------
基于 ```Go``` 语言的 ```discuz!x 3.x To xiunobbs 4.x``` 转换工具

### 开发进度
- 基础构架 ✔
- 数据转换 ✔
- 附件转换
- 数据优化 ✔
> 版主转换

### 编译指南
- 拉取主项目 ```git clone https://github.com/skiy/xiuno-tools.git``` 
- 进入项目目录, 执行 ```go get```
- 编译程序 ```go install``` 或者 ```go build```
- 完成，文件在 ```$GOPATH/bin``` 或在 ```当前目录(go build)``` 下

**温馨提示:**
> 如果已配置好``GOBIN``或者将 ``$GOPATH/bin`` 环境变量，   
即可以在任何目录下执行 **``xiuno-tools``** 启动本程序。   
程序必须有**可执行权限**。   

**工具使用教程**
- 先建一个 xiuno4 论坛。
- 下载本程序（选择运行平台），Linux、MacOS 需要可执行权限。
- 配置```confit.toml```, 执行本程序 ```./xiuno-tools```(Windows 平台下, 建议使用 ```cmd```控制台, 执行```xiuno-tools.exe```)
- 登录后台，记得更新缓存统计。

### 配置文件说明
```toml
[setting]

# 日志配置
[log]
    # 日志等级
    level = "all"
    # 目录
    path = "logs"
    # 是否输出错误位置
    trace = false

# 数据库配置
[database]
    # XiunoBBS
    [[database.xiuno]]
        type = "mysql"      # 数据库类型(不可修改)
        host = "127.0.0.1"  # IP
        port = "3306"       # 端口
        user = "root"       # 数据库用户名
        pass = "123456"     # 密码
        name = "xiuno"      # 数据库名
        prefix = "bbs_"     # 表前缀
        charset = "utf8"    # 字符集
        debug = false     # 日志调试

    # Discuz!X
    [[database.discuz]]
        type = "mysql"
        host = "127.0.0.1"
        port = "3306"
        user = "root"
        pass = "123456"
        name = "dx"
        prefix = "pre_"
        charset = "utf8"
        debug = false

    # UCenter
    [[database.uc]]
        type = "mysql"
        host = "127.0.0.1"
        port = "3306"
        user = "root"
        pass = "123456"
        name = "dx"
        prefix = "pre_ucenter_"
        charset = "utf8"
        debug = false

# 需要转换的表配置
[tables]
    [tables.xiuno]
        # 用户表
        [tables.xiuno.user]
            # 表名
            name = "user"
            # 是否转换
            convert = true
            # 每次更新条数(留空或 < 2, 则默认为 1 条)
            batch = 100

        # 用户组表
        [tables.xiuno.group]
            # 表名
            name = "group"
            # 是否转换
            convert = true
            # 是否使用 xiunobbs 官方用户组
            official = true

        # 版块表
        [tables.xiuno.forum]
            # 表名
            name = "forum"
            # 是否转换
            convert = true

        # 附件表
        [tables.xiuno.attach]
            # 表名
            name = "attach"
            # 是否转换
            convert = true
            # 每次更新条数(留空或 < 2, 则默认为 1 条), 单条导入时, 错误不会导致程序退出
            batch = 1

        # 主题表
        [tables.xiuno.thread]
            # 表名
            name = "thread"
            # 是否转换
            convert = true
            # 每次更新条数(留空或 < 2, 则默认为 1 条)
            batch = 100
            # 取 >= TID 的数据。当上次转换出错时, 记录此 TID, 方便再次导入
            last_tid = 0

        # 帖子表
        [tables.xiuno.post]
            # 表名
            name = "post"
            # 是否转换
            convert = false
            # 每次更新条数(留空或 < 2, 则默认为 1 条)
            batch = 100
            # 取 >= PID 的数据。当上次转换出错时, 记录此 PID, 方便再次导入
            last_pid = 0

        # 置顶帖子表
        [tables.xiuno.thread_top]
            # 表名
            name = "thread_top"
            # 是否转换
            convert = true

        # 我的主题表
        [tables.xiuno.mythread]
            # 表名
            name = "mythread"
            # 是否转换
            convert = true

        # 我的帖子表
        [tables.xiuno.mypost]
            # 表名
            name = "mypost"
            # 是否转换
            convert = true

# 扩展功能
[extension]
    [forum]
        # 是否导入论坛版主 (不建议使用)
        moderators = false

    [file]
        # 是否启用转移附件文件功能
        open = false

        # XiunoBBS 论坛绝对路径
        xiuno_path = ""
        # Discuz!X 论坛绝对路径
        discuzx_path = ""

        # 附件转移
        attach = false
        # 头像转移
        avatar = false
        # 版块 ICON 转移
        icon = false

    [extension.group]
        # 是否启用此功能
        open = true
        # Discuz 游客用户组 ID
        guest_gid = 7
        # 管理员 UID
        admin_id = 1
        # 添加"删除用户的权限"的用户组 gid: 1,2
        delete_user_power = "1,2"

    [extension.user]
        # 是否修正用户主题数和帖子数(帖子数=主题+回复),比较耗时
        total = true
        # 修正 gid 为 101 的用户及用户组
        normal_user = true

    [extension.thread_post]
        # 是否修正主题的 lastpid 和 lastuid,比较耗时
        fix_last = true
        # 是否修正主题内附件统计数量
        thread_attach_total = true
        # 是否修正帖子内附件统计数量
        post_attach_total = true

```

### 注意事项
本程序使用 ***```go mod```*** 标准库，需要 ***```go1.11 +```*** 的开发环境。

### 更新日志

### 使用到的开源项目
- XiunoBBS https://bbs.xiuno.com
- Discuz!X http://www.discuz.net
-

- https://github.com/gogf/gf (基础框架)
- https://github.com/frustra/bbcode (内容 BBCODE 转 HTML)

### 作者
Author: Skiychan   
Email : dev@skiy.net   
Link  : https://www.skiy.net    

### 开源协议
MIT 协议
