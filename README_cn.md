# Meepo
[![Telegram](https://img.shields.io/badge/Telegram-online-brightgreen.svg)](https://t.me/meepoDiscussion)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/PeerXu/meepo/pulls)

Meepo的目标是以便捷的, 去中心化的形式发布服务.

**本项目还处于初期版本, 接口变动会相对频繁, 请留意.**


## 起因

在传统的客户端-服务端架构的网络中, 服务端所在的网络需要能够被客户端访问, 服务端才能正常提供服务.

但是, 由于各种原因, 导致服务端没有足够的资源去暴露端口.

因此, 作者提供了一个工具, 使得客户端可以访问无法暴露端口的服务端.


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

暂时不支持`chocolatey`, 需要从[release](https://github.com/PeerXu/meepo/releases/latest)下载对应版本并手动安装.


## 快速入门

### 访问未暴露公有IP的SSH服务

现有`bob`和`alice`两个节点, `alice`处于防火墙后(无公网IP).

`bob`需要通过`ssh`服务(22端口)连接到`alice`.

1. 在`alice`上, 初始化并且运行`Meepo`服务

```bash
alice$ meepo config init id=alice
alice$ meepo serve
```
**注意: 如果在初始化时为指定ID, 系统会在启动时, 随机分配一个ID.**

通过whoami子命令可以校验`Meepo`服务启动是否成功.

```bash
alice$ meepo whoami
# output:
alice
```

2. 在`bob`上, 运行`Meepo`服务

```bash
bob$ meepo config init id=bob
bob$ meepo serve
```

通过whoami子命令可以校验`Meepo`服务启动是否成功.

```bash
bob$ meepo whoami
# output:
bob
```

3. 在`bob`通过SSH连接到`alice`

```bash
bob$ eval $(meepo ssh bob@alice)
```

等待片刻, `bob`就会连接上`alice`.


### 访问未暴露公有IP的HTTP服务

现有`bob`和`alice`两个节点, `alice`处于防火墙后(无公网IP).

`bob`需要访问`alice`提供的http服务(80端口).

(假定已经按照上节内容配置)

1. 在`bob`上创建连接到`alice`的`http`服务的`teleportation`.

```bash
bob$ meepo teleport -n http -l :8080 alice :80
# output:
Teleport SUCCESS
Enjoy your teleportation with [::]:8080
```

这时候已经成功建立连接, 可以通过 `http://127.0.0.1:8080` 访问`alice`提供的`http`服务.

2. 查看连接情况

```bash
bob$ meepo teleportation list
# output:
+------+-----------+--------+--------------------+--------------------+----------+
| NAME | TRANSPORT | PORTAL |       SOURCE       |        SINK        | CHANNELS |
+------+-----------+--------+--------------------+--------------------+----------+
| http | alice     | source | tcp:[::]:8080      | tcp::80            |        0 |
+------+-----------+--------+--------------------+--------------------+----------+
```

3. 关闭连接

```bash
bob$ meepo teleportation close http
# output:
Teleportation closing
```


## 原理

TBD


## 特性

### 自组网 (Selfmesh)

自组网是`Meepo`的特性, 允许`Meepo`服务提供`Signaling Server`(`WebRTC`建立连接需要交换信令, `Signaling Server`提供交换的服务).

举个简单的例子:

比如现在有三个节点, ID分别为`bob`, `alice`和`eve`.

`bob`与`alice`建立了`Transport`.

`eve`与`alice`建立了`Transport`.

当自组网特性未启用时, 如果需要建立`bob`与`eve`的`Transport`时, 使用的是默认的`Signaling Server`.

未启用自组网时, 交换`Signaling`示意图:

```
bob --- Default Signaling Server --- eve
```

但是, 当自组网特性启用后, 会采用`alice`做为`Signaling Server`, 而不需要使用默认的`Signaling Server`.

启用自组网后, 交换`Signaling`示意图:

```
bob --- alice(Signaling Server) --- eve
```

启用自组网功能十分简单, 只需要将`Meepo`配置的`asSignaling`字段为`true`并重启`Meepo`服务即可.

```bash
# bob:
bob$ meepo config set asSignaling=true
# restart meepo

# alice:
alice$ meepo config set asSignaling=true
# restart meepo

# eve:
eve$ meepo config set asSignaling=true
# restart meepo
```

### Socks5代理

[Socks5](https://zh.wikipedia.org/wiki/SOCKS)是我们常用的网络代理协议之一.

`Meepo`允许用户使用`Socks5`代理访问其他`Meepo`节点提供的服务.

例如节点ID为`hello`, 在80端口上提供了`http`服务.

在完成配置后, 可以直接在浏览器上访问`http://hello.mpo/`访问对应的内容.

域名是采用简单的定义规则, 是`<id>.mpo`.

默认参数下, `Socks5`代理监听地址为`127.0.0.1:12341`.

接下来介绍一下使用方法.

现在有两个节点, ID分别为`bob`和`alice`.

在`alice`节点上, 提供了两个服务, 分别是`ssh`服务(22端口)和`http`服务(80端口).

在`bob`节点上通过`Socks5`代理方便地访问`alice`提供的`ssh`服务和`http`服务.

首先, 需要在`bob`节点上启用`Socks5`功能. (默认情况下已启用)

```bash
bob$ meepo config set proxy.socks5=file://<(echo 'aG9zdDogMTI3LjAuMC4xCnBvcnQ6IDEyMzQxCg=='|base64 -d)
# restart meepo
```

这时候, 已经完成配置工作.

通过`bob`的`Socks5`代理访问`alice`节点的`http`服务.

下面用`curl`举例子.

```bash
bob$ curl -x socks5h://127.0.0.1:12341 http://alice.mpo/
# ...
```

通过`bob`的`Socks5`代理访问`alice`节点的`ssh`服务.

```bash
bob$ ssh -o ProxyCommand='nc -X 5 -x 127.0.0.1:12341 %h %p' bob@alice.mpo
```

Socks5配置请参考各个系统的配置方法.


## 安全

在默认配置下, `Meepo`之间的连接是不需要安全认证的. 这样带来了一定的便捷性, 同时也引入了安全问题.

所以`Meepo`支持共享密钥(secret)的形式增加安全认证机制.

例子:

现在环境中有3个节点, 分别为`alice`, `bob`和`eve`.

假设`alice`和`bob`之间采用共享密钥的形式通信的话, 那么`alice`和`bob`是能够建立连接的.

但是因为`eve`并没有获取到密钥, 所以如果`alice`或`bob`想连接到`eve`上, 是无法成功的, 当然`eve`也无法连接到他们.

```bash
# alice未执行初始化时
alice$ meepo config init id=Alice auth.secret=AliceAndBob
# 或alice已经初始化
alice$ meepo config set auth.secret=AliceAndBob
# ...
alice$ meepo serve
# ...
alice$ meepo whoami
# output:
alice

# 如果Bob已经执行了初始化:
bob$ meepo config set auth.secret=AliceAndBob
# ...
bob$ meepo serve
# ...
bob$ meepo whoami
# Bob
bob$ meepo transport new alice
# wait a few seconds...
bob$ meepo transport list
# output:
+-------+-----------+
| PEER  |   STATE   |
+-------+-----------+
| alice | connected |
+-------+-----------+
```

如果`eve`也需要加入到`alice`和`bob`所组成的网络, 那么需要配置相同的secret. 可以在`alice`或`bob`中导出配置, 并且共享给`eve`.

```bash
# alice:
alice$ meepo config get auth | base64
# output:
bmFtZTogc2VjcmV0CnNlY3JldDogQWxpY2VBbmRCb2IK

# eve:
eve$ meepo config set auth=file://<(echo 'bmFtZTogc2VjcmV0CnNlY3JldDogQWxpY2VBbmRCb2IK' | base64 -d)
# ...
eve$ meepo serve
# ...
eve$ meepo transport new alice
# wait a few seconds...
eve$ meepo transport list
+-------+-----------+
| PEER  |   STATE   |
+-------+-----------+
| alice | connected |
+-------+-----------+
```


## 常见问题

### 修改配置时出现 permission denied

该问题通常出现在操作系统是Linux的情况下.

因为默认配置文件存放在 `/etc/meepo/meepo.yaml`(通过snap安装则在 `/var/snap/meepo/current/etc/meepo.yaml`)下, 所以修改需要对应权限.

```bash
$ sudo bash -c "meepo config set auth=file://<(echo '...'|base64 -d)"
```


## 计划

如果有不错的想法, 不妨通过[Telegram](https://t.me/meepoDiscussion)或[issues](https://github.com/PeerXu/meepo/issues)联系.

- [x] SSH连接端口复用
- [ ] 缩短gather时间
- [ ] 工作原理文档的补全
- [x] 中英文档的补全
- [x] 连接变得可管理
- [x] 支持socks5 proxy
- [ ] 支持http proxy
- [ ] 支持proxy auto-config
- [x] 支持port forward
- [x] 自组网功能
- [x] SignalingEngine认证功能


## 为Meepo做贡献

`Meepo`是一个免费且开源的项目, 欢迎任何人为其开发和进步贡献力量.

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
