# ssgo
A Shadowsocks client with only core features.

***Do One Thing and Do It Well.***
## Features
- SOCKS5 TCP proxy. UDP is **NOT** supported and never will be.
- PAC(Proxy auto-config). Only proxy unreachable host.
- Only AES-256-CFB encryption

## Install
```sh
go get github.com/damoye/ssgo
```
## Usage
```
Usage of ssgo:
  -s string
    server address
  -k string
    password
  -l string
    SOCKS5 server address (default "127.0.0.1:1080")
  -h string
    PAC server address (default "127.0.0.1:8090")
```
### 1. Start ssgo
```sh
ssgo -s [server_address] -k [password]
```
### 2. Config PAC
Config PAC to http://127.0.0.1:8090/proxy.pac
#### OS X:
System Preferences -> Network -> Advanced -> Proxies -> Automatic Proxy Configuration

- Fill the Proxy Configuration File URL with http://127.0.0.1:8090/proxy.pac
- Click OK
- Click Apply

#### Windows:
TODO

## TODO
- Config PAC in Windows
- Test coverage
