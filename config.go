package main

type config struct {
	ServerAddr string `json:"server_addr"`
	Password   string `json:"password"`
	HTTPAddr   string `json:"http_addr"`
	LocalPort  int    `json:"local_port"`
}
