### Xiuno Tools
------

### 开源协议
MIT

### TODO
- ✔XiunoBBS 3.x 升级 XiunoBBS 4.x **[参考](https://gitee.com/xiuno/xiunobbs/blob/master/tool/xn3_to_xn4.php)**
- Discuz!X 3.x 转换 XiunoBBS 4.x [说明](docs/dx3ToXn4/)

### 使用教程
- 拉取依赖库 ```go get -v github.com/skiy/xiuno-tools```   
- 编译程序 ```go install```
- 完成，文件在 ```$GOPATH/bin``` 里

**温馨提示:**
> 如果已配置好``GOBIN``或者将 ``$GOPATH/bin`` 环境变量，   
即可以在任何目录下执行 ``xiuno-tools`` 启动本程序。   
程序必须有**可执行权限**。   

**工具使用教程**
- 先建一个 xiuno4 论坛。
- 下载本程序（选择运行平台），Linux、MacOS 需要可执行权限。
- 登录后台，记得更新缓存统计。

**在 VNC 下后台运行方式执行更新**
> - 在 VNC 窗口1执行命令 ```xiuno-tools > update.log```；
> - 使用 22 端口方式，进入命令行窗口2进入系统进入 ```update.log```所在的目录，并执行 ```tail -f update.log``` 监听日志；
> - 在 VNC 窗口1 按照 日志窗口2的提示输入数据库相关信息，输入完成后，即可关闭日志窗口及 VNC 窗口。
> - 待过一段时间后，再登录系统查看 ```update.log``` 查看转换日志是否完成。

### 更新日志

### 鸣谢
- XiunoBBS https://bbs.xiuno.com
- Discuz http://www.discuz.net
-
- https://github.com/go-sql-driver/mysql
- https://github.com/PuerkitoBio/goquery
- https://github.com/frustra/bbcode
- https://github.com/mozillazg/go-pinyin

### 作者
Author: Skiychan   
Email : dev@skiy.net   
Link  : https://www.skiy.net      