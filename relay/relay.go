package relay

import (
	"io"
	"log"
	"net"

	"github.com/damoye/ssgo/consts"
	"github.com/damoye/ssgo/encrypt"
	"github.com/damoye/ssgo/socks5"
)

func pipe(dst io.Writer, src io.Reader, ch chan error) {
	_, err := io.Copy(dst, src)
	ch <- err
}

func handleConn(c net.Conn, server, password string) {
	defer c.Close()
	target, err := socks5.Handshake(c)
	if err != nil {
		log.Print("handshake: ", err)
		return
	}
	rc, err := net.Dial("tcp", server)
	if err != nil {
		log.Print("dial: ", err)
		return
	}
	defer rc.Close()
	targetStr := target.String()
	rc = encrypt.NewConn(rc, password)
	if _, err = rc.Write(target); err != nil {
		log.Print("write: ", err)
		return
	}
	log.Printf("relay %s - %s", c.RemoteAddr(), targetStr)
	ch := make(chan error, 1)
	go pipe(rc, c, ch)
	go pipe(c, rc, ch)
	if err = <-ch; err != nil {
		log.Print("pipe: ", err)
	}
	log.Printf("relay %s / %s", c.RemoteAddr(), targetStr)
}

// Start starts to relay TCP connection
func Start(server, password string) {
	ln, err := net.Listen("tcp", consts.SOCKS5Addr)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				panic(err)
			}
			go handleConn(conn, server, password)
		}
	}()
}
