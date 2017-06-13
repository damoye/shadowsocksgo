package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/damoye/ssgo/config"
)

func main() {
	conf := config.Config{}
	flag.StringVar(&conf.ServerAddr, "s", "", "server address")
	flag.StringVar(&conf.Password, "k", "", "password")
	flag.StringVar(&conf.LocalAddr, "l", "127.0.0.1:1080", "SOCKS5 server address")
	flag.StringVar(&conf.HTTPAddr, "h", "127.0.0.1:8090", "PAC server address")
	flag.Parse()
	if conf.ServerAddr == "" || conf.Password == "" {
		flag.Usage()
		return
	}
	b, _ := json.Marshal(conf)
	log.Print("Config: ", string(b))
	log.Print("Initializing")
	startHTTP(&conf)
	startTCPRelay(&conf)
	log.Print("Started")
	log.Printf("Please change PAC to http://%s/proxy.pac", conf.HTTPAddr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("Ended")
}
