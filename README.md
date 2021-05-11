# Meepo
[![Telegram](https://img.shields.io/badge/Telegram-online-brightgreen.svg)](https://t.me/meepoDiscussion)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/PeerXu/meepo/pulls)

[Chinese](./README_cn.md)

Meepo aims to publish network service more easy and decentralized.

**This project still in progress**


## Install

### Linux

```bash
$ sudo snap install meepo
```

### macOS

```bash
$ brew install PeerXu/tap/meepo
```

### Windows

Not support `chocolatey` now, install meepo manually from [release](https://github.com/PeerXu/meepo/releases/latest).


## Quick Start

### Access ssh server behind firewall or NAT

There are two nodes, `bob` and `alice`, `alice` behind firewall (without public IP address).

`bob` want to connect to `alice` with ssh service.

1. On `alice`, run `Meepo` service

```bash
alice$ meepo config init id=alice
alice$ meepo serve
```
**NOTE: When initial meepo without ID, Meepo will set a random ID on startup.**

Use `whoami` subcommand to verify `Meepo` service was started or not.

```bash
alice$ meepo whoami
# OUTPUT:
alice
```

2. On `bob`, run `Meepo` service

```bash
bob$ meepo config init id=bob
bob$ meepo serve
```

Use `whoami` subcommand to verify `Meepo` service was started or not.

```bash
bob$ meepo whoami
# OUTPUT:
bob
```

3. On `bob`, connect to `alice` with ssh client.

```bash
bob$ eval $(meepo ssh bob@alice)
# wait a few seconds
# ...
```

### Access http server behind firewall or NAT

There are two nodes, `bob` and `alice`, `alice` behind firewall (without public IP address).

`bob` want to access `http` service which provide by `alice`.

1. On `bob`, new a `teleportation` to access `http` service.

```bash
bob$ meepo teleport -n http -l :8080 alice :80
# OUTPUT:
Teleport SUCCESS
Enjoy your teleportation with [::]:8080
```

Now, enter `http://127.0.0.1:8080` on browser to access `http` service.

2. Teleportation

```bash
bob$ meepo teleportation list
# OUTPUT:
+------+-----------+--------+--------------------+--------------------+----------+
| NAME | TRANSPORT | PORTAL |       SOURCE       |        SINK        | CHANNELS |
+------+-----------+--------+--------------------+--------------------+----------+
| http | alice     | source | tcp:[::]:8080      | tcp::80            |        0 |
+------+-----------+--------+--------------------+--------------------+----------+
```

3. Close teleportation

```bash
bob$ meepo teleportation close http
# OUTPUT:
Teleportation closing
```


## Principle

TBD


## Features

### Selfmesh

Selfmesh, a feature to help `Meepo` nodes to connect each other without `Default Signaling Server` (`WebRTC` need to exchange `signaling` when build connections).

Example:

There are three nodes, `bob`, `alice` and `eve`.

`bob` are built a `transport` with `alice`.

`eve` are built a `transport` with `alice`.

When disable selfmash, if `bob` want to build a `transport` to `eve`, it is using `Default Signaling Server` to exchange `signaling`.

Exchange path when disable selfmash:

```
bob --- Default Signaling Server --- eve
```

When enable selfmash, `alice` will be a `Signaling Server` to exchange `signaling` between `bob` and `eve`.

Exchange path when enable selfmash:

```
bob --- alice(Signaling Server) --- eve
```

It is easy to enable selfmash, set `asSignaling` field to `true` and reboot `Meepo`.

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


### Socks5 Proxy

[Socks5](https://zh.wikipedia.org/wiki/SOCKS) is a usual proxy protocol.

`Meepo` allow user to access service, which provided other `Meepo` node, through `Socks5` proxy.

For example, node ID is `hello`, and node is serving a `http` server(port 80).

We can enter `http://hello.mpo` on browser to access `http` service, when setup socks5 proxy on system and `Meepo`.

The naming rule of domain is `<id>.mpo`.

On default parameters, `Socks5` proxy listen address is `127.0.0.1:12341`.

There are two nodes, `bob` and `alice`.

Two services are running on `alice`, `ssh` service(port 22) and `http` service(port 80).

On `bob`, we can access `ssh` service and `http` service provided by `alice` through `Socks5` proxy on `bob`.

First, we need to enable `Socks5` proxy on `bob`. (enabled by default)

```bash
bob$ meepo config set proxy.socks5=file://<(echo 'aG9zdDogMTI3LjAuMC4xCnBvcnQ6IDEyMzQxCg=='|base64 -d)
# restart meepo
```

Secondly, setup OS proxy setting.

Example:

1. Access `http` service on `bob`

```bash
bob$ curl -x socks5h://127.0.0.1:12341 http://alice.mpo/
# ...
```

2. Access `ssh` service on `bob`

```bash
bob$ ssh -o ProxyCommand='nc -X 5 -x 127.0.0.1:12341 %h %p' bob@alice.mpo
# ...
```


## Security

TBD


## FAQ

TBD


## Roadmap

TBD


## Contributing

`Meepo` is an open source project, welcome every one to contribute codes and documents or else to help `Meepo` to be stronger.

* If any problems about `Meepo`, feel free to open an [issue](https://github.com/PeerXu/meepo/issues).
* If any problems about `Meepo`, feel free to contact us with [Telegram](https://t.me/meepoDiscussion).
* Main branch is used to release stable version, please commit [pull request](https://github.com/PeerXu/meepo/pulls) to dev branch.
* Please feel free to commit bug fix to dev branch.


## Donations

If `Meepo` is helpful for you, welcome to donate to us.

### Telegram

[https://t.me/meepoDiscussion](https://t.me/meepoDiscussion)

### BTC

![BTC](./donations/btc.png)

36PnaXCMCtKLbkzVyfrkudhU6u8vjbfax4

### ETH

![ETH](./donations/eth.png)

0xa4f00EdD5fA66EEC124ab0529cF35a64Ee94BFDE


## Contributer

[PeerXu](https://github.com/PeerXu) (pppeerxu@gmail.com)


## License

MIT
