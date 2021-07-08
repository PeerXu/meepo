# Meepo
[![Telegram](https://img.shields.io/badge/Telegram-online-brightgreen.svg)](https://t.me/meepoDiscussion)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/PeerXu/meepo/pulls)

[Chinese](./README_cn.md)

Meepo aims to publish network service more easy and decentralized.

**This project still in progress**

**BREAKING CHANGE, v0.6 or higher version are not to keep backward compatible.**


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

### Setup

`Meepo` is using [`ED25519 Algorithm`](https://en.wikipedia.org/wiki/EdDSA#Ed25519) as identity algorithm.

Run `meepo serve` to start `Meepo Service`.

```bash
$ meepo serve
```

Run `meepo whoami` to get `MeepoID` of `Meepo Service`.

```bash
$ meepo whoami
# OUTPUT:
61pwmvz1lpm038xwku3njzj21h9na71clie4wv9px1kcxfk49z4
```

Run `meepo shutdown` to shutdown `Meepo Service`.

```bash
$ meepo shutdown
# OUTPUT:
Meepo shutting down
```

Cause we start `Meepo Service` without `Identity File`, `Meepo Service` generate a `Random Identity` to access `Meepo Network`.

We can use `meepo keygen` or `ssh-keygen` to generate `Identity File`.

NOT support `OpenSSH Private Key` with `passphrase` now.

```bash
$ meepo keygen -f meepo.pem
# OR
$ ssh-keygen -t ed25519 -f meepo.pem
```

After generated a `Identity File`, start `Meepo Service` with `Identity File`.

```bash
$ meepo serve -f meepo.pem
```

When `Meepo Service` was started, use `meepo whoami` to get `MeepoID`.

```bash
$ meepo whoami
# OUTPUT:
63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey
```

### Deploy a service to `Meepo Network`

`alice` want to deploy a `HelloWorld Service` to `Meepo Network`.

We make a `HelloWorld Service` now.

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

Start `Meepo Service` and get `MeepoID`.

```bash
# alice:terminal:2
alice$ meepo serve
alice: meepo whoami
# OUTPUT:
63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey
```

Now, we was deployed a `HelloWorld Service` to `Meepo Network`.


### Access deployed `Service` though `Meepo Network`

If `bob` want to access the `HelloWorld Sevice`, deployed by `alice`, `bob` need to start `Meepo Service` too.

But if `bob` do not need to deploy any service to `Meepo Network`, `Random Identity` is good enough.

```bash
# bob:terminal:1
bob$ meepo serve
```

Run `meepo teleport`, to new a `Teleportation` to connect to the `HelloWorld Service` was deployed by `alice`.

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

When `bob` do not need to access the `HelloWorld Service`, run `meepo teleportation close` to close `Teleportation`.

```bash
# bob:terminal:1
bob$ meepo teleportation close alice:http:8080
# OUTPUT:
Teleportation is closing
```


## Principle

TBD


## Features

### Selfmesh

Selfmesh, a feature to help `Meepo Service` to connect each other without `Default Signaling Server` (`WebRTC` need to exchange `signaling` when build connections).

Example:

There are three nodes, `alice`, `bob` and `eve`.

`alice` are built a `transport` with `bob`.

`eve` are built a `transport` with `bob`.

When disable selfmash, if `alice` want to build a `transport` to `eve`, it is using `Default Signaling Server` to exchange `signaling`.

Exchange path when disable selfmash:

```
alice --- Default Signaling Server --- eve
```

When enable selfmash, `bob` will be a `Signaling Server` to exchange `signaling` between `alice` and `eve`.

Exchange path when enable selfmash:

```
alice --- bob(Signaling Server) --- eve
```

`Selfmash` feature was enabled in default.

### `SOCKS5 Proxy`

[SOCKS5](https://zh.wikipedia.org/wiki/SOCKS) is a usual proxy protocol.

`Meepo` allow user to access service, which provided other `Meepo Service`, through `SOCKS5 Proxy`.

For example, `alice` `MeepoID` is `63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey`, and `alice` was deployed a `HelloWorld Service`(port 80).

We can enter `http://63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo` on browser to access the `HelloWorld Service`, when setup `SOCKS5 Proxy` on system and `Meepo`.

The naming rule of domain is `<id>.mpo`.

On default parameters, `SOCKS5 Proxy` listen on `127.0.0.1:12341`.

There are `alice` and `bob`.

Two services are running on `alice`, `SSH Service`(port 22) and `HTTP Service`(port 80).

On `bob`, we can access `SSH Service` and `HTTP Service` provided by `alice` through `SOCKS5 Proxy`.

Example:

1. Access `HTTP Service` on `bob`

```bash
bob$ curl -x socks5h://127.0.0.1:12341 http://63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo/
# ...
```

2. Access `SSH Service` on `bob`

```bash
bob$ ssh -o ProxyCommand='nc -X 5 -x 127.0.0.1:12341 %h %p' bob@63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo
```

## Security

### Authorization

In default parameters, create a `Teleportation` between `Meepo Service` without `authorization`.

Everyone can access the service without `authorization`.

If you do not want anyone can access the service, please setup `authorization` for `Meepo Service`.

Example:

There are `alice` and `bob`.

`alice` `MeepoID` is `63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey`.

`alice` deploy the `HTTP Service`(port 80) and `SSH Service`(port 22), and setup `authorization` with `secret`, `secret` is `AliceAndBob`.

```bash
alice$ cat << EOF > meepo.yaml
meepo:
  auth:
    name: secret
    secret: AliceAndBob
EOF

# Shutdown Meepo Service
alice$ meepo shutdown
# ...

# Start Meepo Service with config file
alice$ meepo servce --config meepo.yaml --identity-file meepo.pem
```

Setup `authorization` is done.

Now, `bob` want to access `HTTP Service` was deployed by `alice`.

`bob` need to add `secret` parameter when `Create Teleportation` or `Teleport`.

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

Sure, `SOCKS Proxy` is support with `authorization`.

`bob` access `HTTP Service` was deployed by `alice` though `SOCKS Proxy`.

```bash
bob$ curl -X socks5h://meepo:AliceAndBob@127.0.0.1:12341 http://63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo
```

`bob` access `SSH Service` was deployed by `alice` though `SOCKS5 Proxy`.

```bash
bob$ ssh -o ProxyCommand='meepo ncat --proxy-type socks5 --proxy 127.0.0.1:12341 --proxy-auth meepo:AliceAndBob %h %p' bob@63eql8p54qpe1jfp1fmuumzge8y6y4ar5uml7nrrf8amqzmutey.mpo
```

### Access Control List

`Meepo` is using `ACL` to control other `Meepo Service` to call `NewTeleportation`.

We can setup `ACL` on config file.

```bash
$ cat meepo.yaml
meepo:
  acl:
    allows:
    - "127.0.0.1:*"
    blocks:
    - "127.0.0.1:22"
```

This acl configuration means we can create `Teleportation` on `127.0.0.1` with any port exclude port 22.

`ACL` configure has two fields, `allows` and `blocks`.

`allows` is a list of `AclPolicy`, which allow matched challenge to create `Teleportation`.

`blocks` is a list of `AclPolicy`, which not allow matched challenge to create `Teleportation`.

`ACL` fllow the rules to run.

1. If challenge triggered `block policies`, then not allow to create `Teleportation`.
2. If challenge triggered `allow policies`, then allow to create `Teleportation`.
3. Not allow to create `Teleportation`.

Let's discuss about `AclPolicy`.

`AclPolicy` format is `source-acl-entity,destination-acl-entity`.

In commons, `source-acl-entity` is `ANY` implicitly if not presents.

`source-acl-entity` and `destination-acl-entity` is `AclEntity`.

`AclEntity` format is `<meepo-id>:<addr-network>:<addr-host>:<addr-port>`.

`addr-network` support `tcp`, `socks5` and `*`.

`addr-host` support `IP Address in IPv4`, `CIDR in IPv4` and `*`.

`addr-port` support network ports and `*`.

Examples:

1. `*` => `*:*:*:*,*:*:*:*`

Match all `Challenge`.

2. `127.0.0.1:22` => `*:*:*:*,*:*:127.0.0.1:22`

Match `Destination.Host` is `127.0.0.1`, `Destination.Port` is `22`.

3. `*:socks5:*:*,*` => `*:socks5:*:*,*:*:*:*`

Match `Source.Network` is `socks5`.

4. `192.168.1.0/24:*` => `*:*:*:*,*:*:192.168.1.0/24:*`

Match `Destination.Host` is `192.168.1.0/24`.

## FAQ

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
