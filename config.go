package main

type config struct {
	ServerAddr string `json:"server_addr"`
	Password   string `json:"password"`
	LocalAddr  string `json:"local_addr"`
	HTTPAddr   string `json:"http_addr"`
}
