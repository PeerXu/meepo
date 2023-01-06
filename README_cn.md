# Meepo
[![Telegram](https://img.shields.io/badge/Telegram-online-brightgreen.svg)](https://t.me/meepoDiscussion)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/PeerXu/meepo/pulls)

Meepo的目标是以便捷的, 去中心化的形式发布服务.

**本项目还处于初期版本, 接口变动会相对频繁, 请留意.**

**由于接口变动, v0.8或更高版本无法向下兼容, 请升级到最新版本.**

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

需要从[release](https://github.com/PeerXu/meepo/releases/latest)下载对应版本并手动安装.

## 快速入门

### 1. 启动`Meepo`实例

```bash
$ meepo serve --no-identity-file --daemon=false
```

### 2. 连接到测试`Meepo`实例

```bash
$ meepo transport new 65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib
```

**注意**

在创建`Transport`时, 获得`Error: transport exist: 65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`是正常现象[1].

### 3. 观察`Transport`连接状态(可选)

```bash
$ meepo transport list
─----------------------------------------------------+-----------+
|                        ADDR                         |   STATE   |
+-----------------------------------------------------+-----------+
| 62vv3lwalqmdb2657f7ax73fem7gkgzmin3w7qyy0sjjfae0f3p | connected |
| 65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib | new       |
+-----------------------------------------------------+-----------+
```

等待片刻, 再执行一遍以上命令, 可以观察到`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`已经是`connected`状态.

```bash
$ meepo transport list
─----------------------------------------------------+-----------+
|                        ADDR                         |   STATE   |
+-----------------------------------------------------+-----------+
| 62vv3lwalqmdb2657f7ax73fem7gkgzmin3w7qyy0sjjfae0f3p | connected |
| 65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib | connected |
+-----------------------------------------------------+-----------+
```

**注意**

如果`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`一直处于`new`, `connecting`等非`connected`状态, 可以试试以下解决方案[2].

### 4. 通过SOCKS5访问`Meepo`实例提供的服务

```bash
$ curl -x socks5h://127.0.0.1:12341 http://65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib.mpo/
Welcome to Meepo Network!
```

我们已经在`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`节点上, 运行着`HTTP`服务, 监听于`127.0.0.1:80`.

上面的`curl`命令, 通过连接到`Meepo`实例提供的`全局SOCKS5`代理服务(默认`127.0.0.1:12341`), 访问`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`的`HTTP`服务.

在这里, `Meepo`提供的`全局SOCKS5`代理服务会把以`.mpo`作为后缀的域名进行解析, 例如:

`http://65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib.mpo/` 将会解析成,

访问`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`节点的, 目的地址为`http://127.0.0.1/`的服务,

意味着, 我们访问的是`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`节点的`HTTP`服务.

### 5. 通过Teleportation(SOCKS5)访问`Meepo`实例提供的服务

```bash
$ meepo teleport --source-network socks5 --listen 127.0.0.1:1087 65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib "*"
teleportation b2554c82 created, listen on :1087

$ curl -x socks5h://127.0.0.1:1087 http://127.0.0.1/
Welcome to Meepo Network!
```

我们通过`teleprot`命令, 创建了一个`Teleportation`.

该`Teleportation`本质上, 是监听着`127.0.0.1:1087`的`SOCKS5`代理服务, 指向`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`节点.

这里, 我们使用`curl`访问时, 通过`SOCKS5`代理后, 访问的地址是`http://127.0.0.1/`.

意味着, 我们访问的是`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`节点的`HTTP`服务.

### 6. 通过Teleportation(Port Forward)访问`Meepo`实例提供的服务

```bash
$ meepo teleport --listen 127.0.0.1:8080 65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib 127.0.0.1:80
teleportation bc244efc created, listen on 127.0.0.1:8080

$ curl http://127.0.0.1:8080/
Welcome to Meepo Network!
```

我们通过`teleport`命令, 创建了一个`Teleportation`.

该`Teleportation`本质上, 是监听着`127.0.0.1:8080`的`TCP`服务, 指向`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`的`127.0.0.1:80`.

所以, 这里我们直接通过`curl`访问`http://127.0.0.1:8080/`,

意味着, 我们访问的是`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`节点的`HTTP`服务.

## 概念

### Transport

TBD

### Teleportation

TBD

### Channel

TBD

## 配置

TBD

## 进阶指南

TBD

## FAQ

TBD
