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

如果发行版不支持`Snapcraft`, 需要从[release](https://github.com/PeerXu/meepo/releases/latest)下载对应版本并手动安装.

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

# open new terminal
$ meepo whoami
62vpszrzczcn7f946ba2ujr0bk311jhlqti357th6346ot2doc8
```

启动`Meepo`实例后, 可以使用`whoami`获取`Meepo`实例的`身份识别地址`, 后面可以提供给其他人访问使用.

### 2. 连接到测试`Meepo`实例

```bash
$ meepo transport new 65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib
```

**注意:** 在创建`Transport`时, 获得`Error: transport exist: 65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`是正常现象[1].

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

**注意:** 如果`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`一直处于`new`, `connecting`等非`connected`状态, 可以试试以下解决方案[2].

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

### 7. 发布服务(可选)

**注意:** 默认情况下, `Meepo`是允许所有`Meepo`实例访问任何地址和端口的, 所以出于安全考虑, 请务必配置好`访问控制列表`后, 再公布`身份识别地址`.

```bash
$ echo -n "<h1>Welcome to Meepo Network!</h1>" > welcome.html
$ python3 -m http.server --bind 127.0.0.1 8080

# open new terminal
$ curl http://127.0.0.1:8080/welcome.html
<h1>Welcome to Meepo Network!</h1>
```

现在, 我们已经启动了一个HTTP服务, 监听着`127.0.0.1:8080`.

由于我们在本机已经运行了`Meepo`实例, 因此我们需要通过另外的`Meepo`实例来访问本机的HTTP服务.

因为`Meepo`已经兼容了`Webassembly`, 所以我们可以在浏览器上运行`Meepo Webassembly Demo`, 来模拟远程访问本机`Meepo`实例发布的服务.

```bash
$ echo -n "http://peerxu.github.io/meepo.html?url=http://127.0.0.1:8080/welcome.html&target=$(meepo whoami)"
```

运行上面的命令后, 会得到Demo的访问地址.

打开网址后, 等待加载完毕, 可以点击`start`按钮.

运行成功后, 按钮内容会变成`running`, 并且在`output`提示栏内显示`meepo started`.

然后点击`newTransport`按钮, 创建与本机`Meepo`实例的`Transport`.

等待一会儿, 点击`listTransports`按钮, 可以看见已经与本机的`Meepo`实例连接上.

这时候, 我们就可以点击`doRequest`按钮, 等待一会儿, 可以在`output`提示栏内看见`Welcome to Meepo Network!`.

## 概念

### Meepo

TBD

### Transport

TBD

### Channel

TBD

### Teleportation

TBD

### Tracker

TBD

## 工作原理

```
Meepo是以WebRTC技术作为主要底层协议, 实现Meepo实例与实例之间的通信功能.
由于, 两个WebRTC节点创建连接前, 需要交换各自的Description, 才能建立连接.
所以Meepo提供了Tracker系统, 使得WebRTC节点在创建连接前, 能够交换Description.

简化的交换Description流程可以理解为, Meepo节点A需要与Meepo节点B创建Transport,
那么节点A作为Tracker-Client, 对Tracker-Server发起请求, 请求将自己的Description(Offer)发送给节点B.
当Tracker-Server能将Offer转发到节点B时, 就会完成这次交换Description的行为.
当Tracker-Server无法将Offer转发到节点B时, 这次请求将会失败,
当节点A记录的所有Tracker-Server都无法完成转发时, 节点A就无法与节点B创建Transport.

接下来我们将会了解到, Tracker-Server是怎么将Offer转发到节点B.
我们将转发行为划分成三种情况讨论.

1. 节点B就是Tracker-Server自身

这种情况下, Tracker-Server也是Meepo节点, 这时候发现Offer的目的身份识别地址是自身,
那么就把相应的请求进行处理, 并且把Description(Answer)返回给节点A, 完成交换Description.

2. Tracker-Server与节点B已经建立Transport

这种情况下, Tracker-Server会把Offer直接转发到节点B上,
节点B收到目的身份识别地址是自身的情况下, 就进行处理并且返回给Tracker-Server,
然后再由Tracker-Server返回给节点A, 完成交换Description.

3. Tracker-Server与节点B没有建立Transport

在这种情况下, 由于Tracker-Server没有与节点B有直接连接,
只能转发给与节点B身份识别地址最接近的已经创建Transport的N个Meepo节点, 希望他们能够处理这个请求.

当执行到上面的情况3时, 我们需要寻找节点B身份识别地址最接近的Meepo节点,
Meepo采用的寻址算法是Kademlia算法的一个变种实现.

Meepo采用ed25519算法的公钥作为身份识别地址. 在使用Tracker系统交换Description时,
会对Description请求进行加密和签名, 保证消息无法被篡改.
只有拥有目的身份识别地址的节点, 才能对消息进行解密.
中途转发的节点都有验证消息的完整性的义务, 与选择是否转发交换Description的请求的权利.
```

## 基础配置

### 1. 身份识别文件

默认情况下, `Meepo`启动时, 会生成随机的`身份识别地址`.

在以下情况下, 有需要使用固定的`身份识别地址`.

* 需要固定`身份识别地址`提供服务.
* 访问的`Meepo`实例, 配置了`白名单`访问策略.

`Meepo`采用的身份识别算法是`ED25519`, 使用`Base36`算法编码`公钥`得到`身份识别地址`.

所以我们可以采用两种方式生成`身份识别文件`.

1. `meepo`命令的`keygen`子命令

```bash
$ meepo keygen
Key generated!
Your identity file has been saved in mpo_id_ed25519
```

2. `ssh-keygen`命令

```bash
$ ssh-keygen -t ed25519 -f mpo_id_ed25519 -P ""
Generating public/private ed25519 key pair.
Your identification has been saved in mpo_id_ed25519
Your public key has been saved in mpo_id_ed25519.pub
The key fingerprint is:
SHA256:u8DHK5AoKAq8wHadWmh9c/7qJHqMS1zbTRDvtIeTgLA ****@****
The key's randomart image is:
+--[ED25519 256]--+
|     .   .       |
|      o . o      |
|     E . o o     |
|          = +    |
|+  .+...S  B .   |
|*+.++*.+ooo o    |
|*.+ o+*oBo .     |
|.. ...o=o+       |
|     oo.+oo.     |
+----[SHA256]-----+
```

**注意:** `Meepo`暂时只支持`passphrase`为空的`身份识别文件`.

可以使用`whoami`子命令读取`身份识别文件`的`身份识别地址`.

```bash
$ meepo whoami -f mpo_id_ed25519
62xz7eexhrgr4lsd7eelyixemh416jom0vg0scj5wyzjc94jfsk
```

启动`Meepo`实例时, 指定`身份识别文件`.

```bash
$ meepo serve --identity-file mpo_id_ed25519
```

### 2. 访问控制列表

`访问控制列表`的配置需要在配置文件中设置, 对应的项是`meepo.acl`.

**注意:** `meepo.acl`的类型是`string`.

默认情况下, `Meepo`实例是允许其他`Meepo`实例访问所有的地址和端口的.

但是出于安全考虑, 我们提供了`访问控制列表`功能, 可以指定允许持有指定`身份识别地址`的`Meepo`实例对指定的目的地址和端口进行访问.

`访问控制列表`的格式如下:

```yaml
- <action1>: "<addr>[,<network>[,<host>:<port>]]"
# ...

# or

- <action2>:
  - "<addr>[,<network>[,<host>:<port>]]"
  # ...
# ...
```

`action1`可以是`allow`或`block`.

`action2`可以是`allows`或`blocks`.

`addr`, `network`, `host`和`port`是支持通配符的.

我们用例子来解析一下`访问控制列表`的规则.

#### 例子

```yaml
- allow: "65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib"
- allow: "*,tcp,127.0.0.1:80"
- block: "*"
```

以上3条规则, 意思分别是:

1. `allow: "65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib"`

允许`65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib`访问任何目的地址和端口.

2. `allow: "*,tcp,127.0.0.1:80"`

允许任何`Meepo`实例访问目的地址为`127.0.0.1`和目的端口为`80`的服务.

3. `block: "*"`

不允许任何`Meepo`实例访问任何地址.

`访问控制列表`的解析规则是顺序执行, 所以如果匹配了任何规则, 那么就会执行相对应的行为.

`network`暂时只支持`tcp`.

## 进阶指南

### 1. 在浏览器中运行`Meepo`实例

[Meepo Webassembly Demo](http://peerxu.github.io/meepo.html?target=65j07gtrxewig4ns5ehlgj21qn15zlphc8e726lqvgrl788zgib)

[About](https://github.com/PeerXu/meepo/tree/main/wasm)

## FAQ

### 1. Error: transport exist

TBD

### 2. Transport一直处于非connected状态

TBD

### 3. Meepo Webassembly Demo点击start后, listTransports列表为空

TBD

## 捐赠

如果`Meepo`对你有帮助, 欢迎对我们进行捐赠.

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
