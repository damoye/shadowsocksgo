# ssgo
A Shadowsocks client with only core features.

***Do One Thing and Do It Well.***
## Features
- SOCKS5 TCP proxy. UDP is **NOT** supported.
- PAC(Proxy auto-config).
- AES-256-CFB encryption **ONLY**.

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
  -h string
        pac http address (default "127.0.0.1:8090")
  -l int
        local socks5 proxy port (default 1080)
```
### Step 0: Start ssgo
```sh
ssgo -s [server_address] -k [password]
```
### Step 1: Config PAC
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
