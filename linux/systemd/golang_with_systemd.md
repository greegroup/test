# 利用systemd部署golang项目的实现方法
 
CentOS 7 继承了 RHEL 7 的新的特性，例如强大的 `systemd`， 而 `systemd` 的使用也使得以往系统服务的 `/etc/init.d` 的启动脚本的方式就此改变， 也大幅提高了系统服务的运行效率。但服务的配置和以往也发生了极大的不同，同时变的简单而易用了许多。

CentOS 7 的服务 `systemctl` 脚本存放在：`/usr/lib/systemd/`，有系统 `system` 和用户 `user` 之分， 即：`/usr/lib/systemd/system` 和 `/usr/lib/systemd/user`

## 配置文件

这里先要说明一下 `unit` 的文件位置，一般主要有三个目录：

```
/lib/systemd/system
/run/systemd/system
/etc/systemd/system
``` 

这三个目录的配置文件优先级依次从低到高，如果同一选项三个地方都配置了，优先级高的会覆盖优先级低的。 系统安装时，默认会将 `unit` 文件放在 `/lib/systemd/system` 目录。如果想要修改系统默认的配置，比如 `nginx.service`，一般有两种方法：

- 在 `/etc/systemd/system` 目录下创建 `nginx.service` 文件，里面写上我们自己的配置。

- 在 `/etc/systemd/system` 下面创建 `nginx.service.d` 目录，在这个目录里面新建任何以`.conf` 结尾的文件，然后写入我们自己的配置。推荐这种做法。

`/run/systemd/system` 这个目录一般是进程在运行时动态创建 unit 文件的目录，一般很少修改，除非是修改程序运行时的一些参数时，即 Session 级别的，才在这里做修改。

### 服务配置
每一个服务以`.service` 结尾，一般会分为 3 部分：`[Unit]`、`[Service]` 和 `[Install]`，就以 nginx 为例吧，具体内容如下：
```
[Unit]
 
Description=nginx - high performance web server
 
Documentation=http://nginx.org/en/docs/
 
After=network.target remote-fs.target nss-lookup.target
 
[Service]
 
Type=forking
 
PIDFile=/usr/local/nginx/logs/nginx.pid
 
ExecStartPre=/usr/local/nginx/sbin/nginx -t -c /usr/local/nginx/conf/nginx.conf
 
ExecStart=/usr/local/nginx/sbin/nginx -c /usr/local/nginx/conf/nginx.conf
 
ExecReload=/bin/kill -s HUP $MAINPID
 
ExecStop=/bin/kill -s QUIT $MAINPID
 
PrivateTmp=true
 
[Install]
 
WantedBy=multi-user.target
```
### 配置项说明
下面分别解释下着三部分的含义

#### [Unit]
- Description : 服务的简单描述

- Documentation ： 服务文档

- After= : 依赖，仅当依赖的服务启动之后再启动自定义的服务单元

#### [Service]
- Type : 启动类型 simple、forking、oneshot、notify、dbus

- Type=simple（默认值）：systemd 认为该服务将立即启动，服务进程不会 fork。如果该服务要启动其他服务，不要使用此类型启动，除非该服务是 socket 激活型
- Type=forking：systemd 认为当该服务进程 fork，且父进程退出后服务启动成功。对于常规的守护进程（daemon），除非你确定此启动方式无法满足需求， 使用此类型启动即可。使用此启动类型应同时指定 PIDFile=，以便 systemd 能够跟踪服务的主进程。
- Type=oneshot：这一选项适用于只执行一项任务、随后立即退出的服务。可能需要同时设置 RemainAfterExit=yes 使得 systemd 在服务进程退出之后仍然认为服务处于激活状态。
- Type=notify：与 Type=simple 相同，但约定服务会在就绪后向 systemd 发送一个信号，这一通知的实现由 libsystemd-daemon.so 提供
- Type=dbus：若以此方式启动，当指定的 BusName 出现在 DBus 系统总线上时，systemd 认为服务就绪。
- PIDFile ： pid 文件路径

- Environment ： 环境变量（可以添加多个）eg ：Environment=REPO_REF=dev

- ExecStartPre ：启动前要做什么，上文中是测试配置文件 －t

- ExecStart：启动

- ExecReload：重载

- ExecStop：停止

- PrivateTmp：True 表示给服务分配独立的临时空间

#### [Install]
- WantedBy：服务安装的用户模式，从字面上看，就是想要使用这个服务的有是谁？上文中使用的是：multi-user.target ，就是指想要使用这个服务的目录是多用户。

## 操作示例

每一个.target 实际上是链接到我们单位文件的集合，当我们执行

```
systemctl enable nginx.service
```

就会在 /etc/systemd/system/multi-user.target.wants/ 目录下新建一个 /usr/lib/systemd/system/nginx.service 文件的链接。

**常用的 service 操作**
```
# 自启动
systemctl enable nginx.service
 
# 禁止自启动
systemctl disable nginx.service
 
# 启动服务
systemctl start nginx.service
 
# 停止服务
systemctl stop nginx.service
 
# 重启服务
systemctl restart nginx.service
 
# 查看Unit定义文件
systemctl cat nginx.service
 
# 编辑Unit定义文件
systemctl edit nginx.service
 
# 重新加载Unit定义文件
systemctl reload nginx.service
 
# 列出已启动的所有unit，就是已经被加载到内存中
systemctl list-units
 
# 列出系统已经安装的所有unit，包括那些没有被加载到内存中的unit
systemctl list-unit-files
 
# 查看服务的日志
journalctl -u nginx.service # 还可以配合`-b`一起使用，只查看自本次系统启动以来的日志
 
# 查看所有target下的unit
systemctl list-unit-files --type=target
 
# 查看默认target，即默认的运行级别。对应于旧的`runlevel`命令
systemctl get-default
 
# 设置默认的target
systemctl set-default multi-user.target
 
# 查看某一target下的unit
systemctl list-dependencies multi-user.target
 
# 切换target，不属于新target的unit都会被停止
systemctl isolate multi-user.target
systemctl poweroff # 关机
systemctl reboot  # 重启
systemctl rescue # 进入rescue模式

```
## Systemd + Golang Demo

下面例子通过 Systemd 部署一个最简单的 Golang Web Server
```
package main
 
import (
 "flag"
 "fmt"
 "log"
 "net/http"
)
 
func main() {
 var address string
 flag.StringVar(&address, "address", "5353", "listen address")
 flag.Parse()
 
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, This is a test for systemd !\n")
 }) // 设置访问路由
 
 log.Printf("Starting server on %s\n", address)
 log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", address), nil))
}
```
编译代码，并将可执行文件同步到远程服务器上
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o systemd-test main.go
rsync -zP ./systemd-test root@gp-aliyun:/usr/local/bin/
```

在远程服务器上编写服务配置，放在 `/etc/systemd/system/` 中
```
[Unit]
Description=Systemd Test
After=network.target
 
[Service]
User=nobody
# Execute `systemctl daemon-reload` after ExecStart= is changed.
ExecStart=/usr/local/bin/systemd-test -address "6060"
[Install]
WantedBy=multi-user.target
```


通过 `systemctl` 启动服务
```
# 每一次修改ExecStart都需执行
systemctl daemon-reload
# 启动
systemctl start systemd-test.service
# 查看状态
systemctl status systemd-test.service
```

状态如下
![](https://github.com/greegroup/test/blob/master/linux/systemd/assets/systemd_01.png)	

可以通过 curl 进行测试：
![](https://github.com/greegroup/test/blob/master/linux/systemd/assets/systemd_02.png)