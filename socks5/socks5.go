package socks5

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	ver5         = 0x05
	methodNoAuth = 0x00
	cmdConnect   = 0x01
	atypIPV4     = 0x01
	atypIPV6     = 0x04
	atypHOST     = 0x03
	repSucceeded = 0x00
)

// Address ...
type Address []byte

func (a Address) String() string {
	var host string
	var port uint16
	switch a[0] {
	case atypIPV4:
		host = net.IP(a[1 : 1+net.IPv4len]).String()
		port = binary.BigEndian.Uint16(a[1+net.IPv4len:])
	case atypIPV6:
		host = net.IP(a[1 : 1+net.IPv6len]).String()
		port = binary.BigEndian.Uint16(a[1+net.IPv6len:])
	case atypHOST:
		host = string(a[2 : 2+a[1]])
		port = binary.BigEndian.Uint16(a[2+a[1]:])
	default:
		return ""
	}
	return net.JoinHostPort(host, strconv.Itoa(int(port)))
}

// Handshake for SOCKS5
func Handshake(conn net.Conn) (Address, error) {
	// read VER, NMETHODS, METHODS
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	if buf[0] != ver5 { // VER: 5
		return nil, errors.New(fmt.Sprint("not socks5:", buf[0]))
	}
	buf = make([]byte, buf[1])
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	noAuthExist := false
	for _, method := range buf {
		if method == methodNoAuth {
			noAuthExist = true
			break
		}
	}
	if !noAuthExist {
		return nil, errors.New(fmt.Sprint("no method noAuth: ", buf))
	}

	// write VER, METHOD
	if _, err := conn.Write([]byte{ver5, methodNoAuth}); err != nil {
		return nil, err
	}

	// read VER, CMD, RSV, ATYP, DST.ADDR, DST.PORT
	buf = make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	if buf[1] != cmdConnect { // CMD: CONNECT
		return nil, errors.New(fmt.Sprint("not connect cmd: ", buf[1]))
	}
	addrType := buf[3]
	switch addrType {
	case atypIPV4:
		buf = make([]byte, 4+2)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
	case atypIPV6:
		buf = make([]byte, 16+2)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
	case atypHOST:
		buf = make([]byte, 1)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
		len := buf[0]
		buf = make([]byte, 1+len+2)
		buf[0] = len
		if _, err := io.ReadFull(conn, buf[1:]); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprint("unknown address type", buf[3]))
	}

	// write VER REP RSV ATYP BND.ADDR BND.PORT
	reply := []byte{ver5, repSucceeded, 0x00, atypIPV4, 0x00, 0x00, 0x00, 0x00, 0x10, 0x10}
	if _, err := conn.Write(reply); err != nil {
		return nil, err
	}

	// return ATYP DST.ADDR DST.PORT
	return Address(append([]byte{addrType}, buf...)), nil
}
