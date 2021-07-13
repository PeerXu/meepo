# Meepo
[![Telegram](https://img.shields.io/badge/Telegram-online-brightgreen.svg)](https://t.me/meepoDiscussion)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/PeerXu/meepo/pulls)

Meepo的目标是以便捷的, 去中心化的形式发布服务.

**本项目还处于初期版本, 接口变动会相对频繁, 请留意.**

**由于接口变动, v0.6或更高版本无法向下兼容, 请升级到最新版本.**


## 安装

### Linux

```bash
$ sudo snap install meepo
```

如果发行版不支持`snap`, 需要从[release](https://github.com/PeerXu/meepo/releases/latest)下载对应版本并手动安装.

### macOS

```bash
$ brew install PeerXu/tap/meepo
```

### Windows

暂时不支持从`chocolatey`安装, 需要从[release](https://github.com/PeerXu/meepo/releases/latest)下载对应版本并手动安装.


## 快速入门

### 初始化

`Meepo`采用[`ED25519算法`](https://en.wikipedia.org/wiki/EdDSA#Ed25519)作为身份标识算法.

运行`meepo serve`命令, 启动`Meepo服务`.

```bash
$ meepo serve
```

运行`meepo whoami`命令, 可以获得`Meepo服务`的`MeepoID`.

```bash
$ meepo whoami
# OUTPUT:
61pwmvz1lpm038xwku3njzj21h9na71clie4wv9px1kcxfk49z4
```

运行`meepo shutdown`命令, 可以关闭`Meepo服务`.

```bash
$ meepo shutdown
# OUTPUT:
Meepo shutting down
```

由于未指定`身份标识文件`, `meepo serve`运行时, 会生成`随机身份`, 方便将要启动的`Meepo服务`能够正常接入`Meepo网络`.

可以通过`meepo keygen`或`ssh-keygen`生成`身份标识文件`.

采用`ssh-keygen`生成`身份识别文件`会同时产生私钥和公钥文件, 忽略公钥文件即可.

暂时不支持带`passphrase`的`OpenSSH私钥`.

```bash
$ meepo keygen -f meepo.pem
# OR
$ ssh-keygen -t ed25519 -f meepo.pem
```

`生成身份识别文件`成功之后, 启动`Meepo服务`指定`身份标识文件`.

```bash
$ meepo serve -i meepo.pem
```

`Meepo服务`启动成功后, 再通过`meepo whoami`获取当前的`MeepoID`.

```bash
$ meepo whoami
# OUTPUT:
63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey
```

### 发布一个服务到`Meepo网络`

假如, 现在`alice`需要发布一个`HelloWorld服务`到`Meepo网络`.

先来实现一个简单的`HelloWorld服务`.

```bash
# alice:terminal:1
alice$ cat << EOF > index.html
<h1>Hello World!</h1>
EOF
alice$ cat index.html
# OUTPUT:
<h1>Hello World!</h1>

alice$ python3 -m http.server 8080

# alice:terminal:2
alice$ curl http://127.0.0.1:8080
# OUTPUT:
<h1>Hello World!</h1>
```

这样我们就为`alice`实现了一个`HelloWorld服务`.

接着启动`Meepo服务`和获取`MeepoID`.

```bash
# alice:terminal:2
alice$ meepo serve
alice: meepo whoami
# OUTPUT:
63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey
```

这时候, 已经完成整个发布流程.


### 访问已经发布到`Meepo网络`的服务

接下来, 我们将讲解`bob`如何访问`alice`的`HelloWorld服务`.

首先, `bob`也需要运行一个`Meepo服务`, 但是由于`bob`并不需要发布服务, 所以使用`随机身份`即可.

```bash
# bob:terminal:1
bob$ meepo serve
```

然后, 运行`meepo teleport`命令, 使`bob`的`Meepo服务`生成一条`Teleportation`通往`alice`的`HelloWorld服务`.

```bash
# bob:terminal:1
bob$ meepo teleport -n alice:http:8080 -l 127.0.0.1:8080 63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey 127.0.0.1:8080
# Wait a few minutes...
# OUTPUT:
Teleport SUCCESS
Enjoy your teleportation with 127.0.0.1:8080

bob$ meepo teleportation list
# OUTPUT:
+-----------------+-----------------------------------------------------+--------+--------------------+--------------------+----------+
|      NAME       |                      TRANSPORT                      | PORTAL |       SOURCE       |        SINK        | CHANNELS |
+-----------------+-----------------------------------------------------+--------+--------------------+--------------------+----------+
| alice:http:8080 | 63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey | source | tcp:127.0.0.1:8080 | tcp:127.0.0.1:8080 |        0 |
+-----------------+-----------------------------------------------------+--------+--------------------+--------------------+----------+
bob$ curl http://127.0.0.1:8080
# OUTPUT:
<h1>Hello World!</h1>
```

当`bob`不再需要访问该服务时, 可以通过`meepo teleportation close`命令关闭`Teleportation`.

```bash
# bob:terminal:1
bob$ meepo teleportation close alice:http:8080
# OUTPUT:
Teleportation is closing
```


## 原理

TBD


## 特性

### 自组网 (Selfmesh)

自组网是`Meepo`的特性, 允许`Meepo服务`成为`Signaling Server`(`WebRTC`建立连接需要交换信令, `Signaling Server`提供交换的服务).

举个简单的例子:

比如现在有三个节点, 分别为`alice`, `bob`和`eve`.

`alice`与`bob`建立了`Transport`.

`eve`与`bob`建立了`Transport`.

当自组网特性未启用时, 如果需要为`alice`和`eve`的创建`Transport`时, 使用的是默认的`Signaling Server`.

未启用自组网时, 交换`Signaling`示意图:

```
alice --- Default Signaling Server --- eve
```

但是, 当自组网特性启用后, 会采用`bob`做为`Signaling Server`, 而不需要使用默认的`Signaling Server`.

启用自组网后, 交换`Signaling`示意图:

```
alice --- bob(Signaling Server) --- eve
```

默认参数下, 已经启动自组网功能.

### SOCKS5代理

[SOCKS5](https://zh.wikipedia.org/wiki/SOCKS)是我们常用的网络代理协议之一.

`Meepo`允许用户使用`SOCKS5`代理访问由`Meepo服务`发布的服务.

假如`alice`的`MeepoID`为`63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey`, 并且发布了`HTTP服务`, 端口为80.

在完成配置后, 可以直接在浏览器上访问`http://63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo/`访问`alice`发布的`HTTP服务`.

域名是采用简单的定义规则, 是`<id>.mpo`.

默认参数下, `SOCKS5`代理监听地址为`127.0.0.1:12341`.

接下来介绍一下使用方法.

现在有两个节点, 分别为`alice`和`bob`.

在`alice`上, 发布两个服务, 分别是`SSH服务`(22端口)和`HTTP服务`(80端口).

在`bob`上, 通过`SOCKS5代理`访问`alice`发布的`SSH服务`和`HTTP服务`.

下面用`curl`举例子.

```bash
bob$ curl -x socks5h://127.0.0.1:12341 http://63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo/
# ...
```

通过`bob`的`Socks5`代理访问`alice`节点的`ssh`服务.

```bash
bob$ ssh -o ProxyCommand='nc -X 5 -x 127.0.0.1:12341 %h %p' bob@63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo
```

SOCKS5配置请参考各个系统的配置方法.


## 安全

### 认证

在默认配置下, `Meepo服务`之间创建`Teleportation`是不需要认证的. 这样带来了一定的便捷性, 同时也引入了安全问题.

所以`Meepo`支持以密钥(`secret`)的形式增加安全认证机制.

例子:

假如存在`alice`和`bob`.

`alice`的`MeepoID`为`63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey`.

`alice`发布`HTTP服务`(80端口)和`SSH服务`(22端口), 并且配置`secret`为`AliceAndBob`.

```bash
alice$ cat << EOF > meepo.yaml
meepo:
  auth:
    name: secret
    secret: AliceAndBob
EOF

alice$ meepo shutdown
# ...

alice$ meepo serve --config meepo.yaml --identity-file meepo.pem
```

配置`secret`的工作已经完成.

这时候, `bob`需要访问`alice`发布的`HTTP服务`.

那么`bob`在创建`Teleportation`时, 需要带上参数`--secret`, 指定`secret`.

```bash
bob$ meepo teleport -n alice-http-80 -s AliceAndBob -l 127.0.0.1:8080 63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey 127.0.0.1:80
# wait a few minutes
# OUTPUT:
Teleport SUCCESS
Enjoy your teleportation with 127.0.0.1:8080

bob$ meepo teleportation list
# OUTPUT:
+---------------+-----------------------------------------------------+--------+--------------------+------------------+----------+
|     NAME      |                      TRANSPORT                      | PORTAL |       SOURCE       |       SINK       | CHANNELS |
+---------------+-----------------------------------------------------+--------+--------------------+------------------+----------+
| alice-http-80 | 63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey | source | tcp:127.0.0.1:8080 | tcp:127.0.0.1:80 |        0 |
+---------------+-----------------------------------------------------+--------+--------------------+------------------+----------+

bob$ curl http://127.0.0.1:8080/
# ...
```

当然, `SOCKS5代理`也是支持`secret`的.

`bob`通过`SOCKS5代理`访问`alice`发布的`HTTP服务`.

```bash
bob$ curl -X socks5h://meepo:AliceAndBob@127.0.0.1:12341 http://63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo
```

`bob`通过`SOCKS5代理`访问`alice`发布的`SSH服务`.

```bash
bob$ ssh -o ProxyCommand='meepo ncat --proxy-type socks5 --proxy 127.0.0.1:12341 --proxy-auth meepo:AliceAndBob %h %p' bob@63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo
```

### 访问控制列表

[`ACL`](https://en.wikipedia.org/wiki/Access-control_list)是一种常用的控制访问权限的手段.

`Meepo`使用`ACL`控制其他`Meepo Service`调用`NewTeleportation`的权限.

配置的`meepo.acl`项可以控制`ACL`的行为.

```bash
$ cat meepo.yaml
meepo:
  acl:
    allows:
    - "127.0.0.1:*"
    blocks:
    - "127.0.0.1:22"
```

例如上面这个配置, 意思是除了端口22之外, 可以在`127.0.0.1`上面任意端口创建`Teleportation`.

`ACL`行为由两个列表影响, `allows`和`blocks`.

`allows`是允许通过的规则(`AclPolicy`).

`blocks`是不允许通过的规则.

`ACL`运行逻辑顺序如下:

1. 如果触发`blocks`规则, 则不允许创建`Teleportation`.
2. 如果触发`allows`规则, 则允许创建`Teleportation`.
3. 不允许创建`Teleportation`.

下面我们来讨论一下规则.

`AclPolicy`的格式是`source-acl-entity,destination-acl-entity`.

通常情况下, `source-acl-entity`是可以省略的.

`source-acl-entity`和`destination-acl-entity`都是`AclEntity`.

`AclEntity`的格式是`<meepo-id>:<addr-network>:<addr-host>:<addr-port>`.

`addr-network`暂时只支持`tcp`, `socks5`和`*`. 

通常情况下`source-acl-entity.addr-network`可选`tcp`, `socks5`和`*`, `destination-acl-entity.addr-network`可选`tcp`和`*`.

`addr-host`暂时只支持`IPv4格式的IP和CIDR`和`*`.

`addr-port`支持正常网络支持的端口和`*`.

例子:

1. `*` => `*:*:*:*,*:*:*:*`

匹配所有`Challenge`.

2. `127.0.0.1:22` => `*:*:*:*,*:*:127.0.0.1:22`

匹配`Destination.Host`为`127.0.0.1`, `Destination.Port`为`22`.

3. `*:socks5:*:*,*` => `*:socks5:*:*,*:*:*:*`

匹配`Source.Network`为`socks5`.

4. `192.168.1.0/24:*` => `*:*:*:*,*:*:192.168.1.0/24:*`

匹配`Destination.Host`为`192.168.1.0/24`.

## 常见问题

### `Transport`无法创建

由于`Transport`是采用`WebRTC协议`, 所以在创建`Transport`会受到[`WebRTC`条件限制](https://webrtcforthecurious.com/zh/docs/03-connecting/#nat%e6%98%a0%e5%b0%84).

所以暂时有些网络情况是无法正常使用`Meepo服务`. 后期会提供其他解决方案来解决这个问题.

### `Windows平台`下, `Daemon模式`无法正常工作

`Windows平台`暂时不支持`Daemon模式`. `--daemon`参数会被忽略.


## 为Meepo做贡献

`Meepo`是一个免费且开源的项目, 欢迎任何人为其开发和进步贡献力量.

* 如果有不错的想法, 不妨通过[Telegram](https://t.me/meepoDiscussion)或[issues](https://github.com/PeerXu/meepo/issues)联系.
* 在使用过程中出现任何问题, 可以通过[issues](https://github.com/PeerXu/meepo/issues)来反馈.
* 在使用过程中出现任何问题, 也可以通过[Telegram](https://t.me/meepoDiscussion)来沟通使用心得.
* 如果还有其他方面的问题与合作, 欢迎联系 pppeerxu@gmail.com .


### 代码提交

* main分支仅作用于稳定版本的发布, [PRs](https://github.com/PeerXu/meepo/pulls)请提交到dev分支.
* Bug修复可以直接提交PR到dev分支.
* 如果有新增功能的想法, 可以先到[issues](https://github.com/PeerXu/meepo/issues)描述想法与对应的实现, 然后fork修改, 最后提交PR到dev分支进行合并.


## 捐赠

如果觉得Meepo能够帮助到你, 欢迎提供适当的捐助来维持项目的长期发展.

### Telegram

[https://t.me/meepoDiscussion](https://t.me/meepoDiscussion)

### BTC

![BTC](./donations/btc.png)

36PnaXCMCtKLbkzVyfrkudhU6u8vjbfax4

### ETH

![ETH](./donations/eth.png)

0xa4f00EdD5fA66EEC124ab0529cF35a64Ee94BFDE


## 贡献者

[PeerXu](https://github.com/PeerXu) (pppeerxu@gmail.com)


## License

MIT
