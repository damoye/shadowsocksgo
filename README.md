# ssgo
There are too many features in Shadowsocks. I just want a TCP proxy with the AES-256-CFB encryption.

***"Do One Thing and Do It Well."***

## Install

```sh
go get github.com/damoye/ssgo
```

## Usage
Start a client connecting to a Shadowsocks server. The client listens on port 1080 for incoming SOCKS5 connections.

```sh
ssgo -s [server_address] -k [password]
```

## TODO

- PAC
- UDP
- Server
- Test coverage
