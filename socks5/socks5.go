package socks5

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

// Address ...
type Address []byte

func (a Address) String() string {
	var host, port string
	switch a[0] {
	case 0x01: // ATYP: IP V4 address
		host = net.IP(a[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(a[1+net.IPv4len]) << 8) | int(a[1+net.IPv4len+1]))
	case 0x04: // ATYP: IP V6 address
		host = net.IP(a[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(a[1+net.IPv6len]) << 8) | int(a[1+net.IPv6len+1]))
	case 0x03: // ATYP: DOMAINNAME
		host = string(a[2 : 2+a[1]])
		port = strconv.Itoa((int(a[2+a[1]]) << 8) | int(a[2+a[1]+1]))
	default:
		return ""
	}
	return net.JoinHostPort(host, port)
}

// Handshake for SOCKS5
func Handshake(conn net.Conn) (Address, error) {
	// read VER, NMETHODS, METHODS
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	if buf[0] != 0x05 { // VER: 5
		return nil, errors.New(fmt.Sprint("not socks5:", buf[0]))
	}
	buf = make([]byte, buf[1])
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	noAuthExist := false
	for _, method := range buf {
		if method == 0x00 { // METHOD: NO AUTHENTICATION REQUIRED
			noAuthExist = true
			break
		}
	}
	if !noAuthExist {
		return nil, errors.New(fmt.Sprint("no method noAuth: ", buf))
	}

	// write VER, METHOD
	if _, err := conn.Write([]byte{0x05, 0x00}); err != nil {
		return nil, err
	}

	// read VER, CMD, RSV, ATYP, DST.ADDR, DST.PORT
	buf = make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	if buf[1] != 0x01 { // CMD: CONNECT
		return nil, errors.New(fmt.Sprint("not connect cmd: ", buf[1]))
	}
	addrType := buf[3]
	switch addrType {
	case 0x01: // ATYP: IP V4 address
		buf = make([]byte, 4+2)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
	case 0x04: // ATYP: IP V6 address
		buf = make([]byte, 16+2)
		if _, err := io.ReadFull(conn, buf); err != nil {
			return nil, err
		}
	case 0x03: // ATYP: DOMAINNAME
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
	if _, err := conn.Write([]byte("\x05\x00\x00\x01\x00\x00\x00\x00\x10\x10")); err != nil {
		return nil, err
	}

	// return ATYP DST.ADDR DST.PORT
	return Address(append([]byte{addrType}, buf...)), nil
}
