package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/damoye/ssgo/consts"
	"github.com/damoye/ssgo/pac"
	"github.com/damoye/ssgo/relay"
)

func main() {
	server := flag.String("s", "", "server address")
	password := flag.String("k", "", "password")
	flag.Parse()
	if *server == "" || *password == "" {
		flag.Usage()
		return
	}
	pac.Start()
	relay.Start(*server, *password)
	log.Print("SOCKS5 is listening at ", consts.SOCKS5Addr)
	log.Print("PAC URL is http://127.0.0.1", consts.HTTPAddr, "/proxy.pac")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
