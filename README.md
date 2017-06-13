# ssgo
A Shadowsocks client only supporting TCP relay and AES-256-CFB encryption.

***Do One Thing and Do It Well.***
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
## TODO
- Test coverage
