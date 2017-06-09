package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"

	"github.com/damoye/shadowsocksgo/encrypt"
	"github.com/damoye/shadowsocksgo/socks5"
)

var (
	infoLog  = log.New(os.Stderr, "INFO: ", log.LstdFlags|log.Lshortfile)
	errorLog = log.New(os.Stderr, "ERRO: ", log.LstdFlags|log.Lshortfile)

	server   = flag.String("s", "", "server address")
	local    = flag.String("l", ":1080", "local address")
	password = flag.String("k", "", "password")
)

func pipe(dst net.Conn, src net.Conn, ch chan error) {
	_, err := io.Copy(dst, src)
	ch <- err
}

func handleConn(c net.Conn) {
	defer c.Close()
	target, err := socks5.Handshake(c)
	if err != nil {
		errorLog.Println("handshake:", err)
		return
	}
	rc, err := net.Dial("tcp", *server)
	if err != nil {
		errorLog.Println("dial:", err)
		return
	}
	defer rc.Close()
	infoLog.Printf("proxy %s <-> %s <-> %s", c.RemoteAddr(), *server, target)
	rc, err = encrypt.NewEncryptedConn(rc, *password, target)
	if err != nil {
		errorLog.Println("newEncryptedConn:", err)
		return
	}
	if _, err = rc.Write(target); err != nil {
		errorLog.Println("write:", err)
		return
	}
	ch := make(chan error, 1)
	go pipe(rc, c, ch)
	go pipe(c, rc, ch)
	if err = <-ch; err != nil {
		errorLog.Println("pipe:", err)
	}
}

func main() {
	flag.Parse()
	if *server == "" || *local == "" || *password == "" {
		flag.Usage()
		return
	}
	ln, err := net.Listen("tcp", *local)
	if err != nil {
		errorLog.Fatalln("listen:", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			errorLog.Println("accept:", err)
			continue
		}
		go handleConn(conn)
	}
}
