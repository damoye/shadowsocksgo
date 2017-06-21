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
	socks5Version  = 5
	socks5AuthNone = 0
	socks5Connect  = 1

	socks5IP4    = 1
	socks5Domain = 3
	socks5IP6    = 4
)

// Addr represents SOCKS address RFC 1928
type Addr []byte

func (a Addr) String() string {
	var host string
	var port uint16
	switch a[0] {
	case socks5IP4:
		host = net.IP(a[1 : 1+net.IPv4len]).String()
		port = binary.BigEndian.Uint16(a[1+net.IPv4len:])
	case socks5IP6:
		host = net.IP(a[1 : 1+net.IPv6len]).String()
		port = binary.BigEndian.Uint16(a[1+net.IPv6len:])
	case socks5Domain:
		host = string(a[2 : 2+a[1]])
		port = binary.BigEndian.Uint16(a[2+a[1]:])
	default:
		return ""
	}
	return net.JoinHostPort(host, strconv.Itoa(int(port)))
}

func readAddr(r io.Reader) (Addr, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	atyp := buf[0]
	switch atyp {
	case socks5IP4:
		buf = make([]byte, 1+net.IPv4len+2)
		buf[0] = atyp
		if _, err := io.ReadFull(r, buf[1:]); err != nil {
			return nil, err
		}
	case socks5IP6:
		buf = make([]byte, 1+net.IPv6len+2)
		buf[0] = atyp
		if _, err := io.ReadFull(r, buf[1:]); err != nil {
			return nil, err
		}
	case socks5Domain:
		buf = make([]byte, 1)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		len := buf[0]
		buf = make([]byte, 1+1+len+2)
		buf[0], buf[1] = atyp, len
		if _, err := io.ReadFull(r, buf[1+1:]); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprint("unknown address type: ", atyp))
	}
	return buf, nil
}

// Handshake for SOCKS5
func Handshake(conn io.ReadWriter) (Addr, error) {
	// +----+----------+----------+
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	if buf[0] != socks5Version {
		return nil, errors.New(fmt.Sprint("not socks5: ", buf[0]))
	}
	buf = make([]byte, buf[1])
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+
	if _, err := conn.Write([]byte{socks5Version, socks5AuthNone}); err != nil {
		return nil, err
	}
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	buf = make([]byte, 3)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	if buf[1] != socks5Connect {
		conn.Write([]byte{socks5Version, 7, 0, 0, 0, 0, 0, 0, 0, 0})
		return nil, errors.New(fmt.Sprint("not connect cmd: ", buf[1]))
	}
	dest, err := readAddr(conn)
	if err != nil {
		return nil, err
	}
	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	reply := []byte{socks5Version, 0, 0, socks5IP4, 0, 0, 0, 0, 0, 0}
	_, err = conn.Write(reply)
	return dest, err
}
