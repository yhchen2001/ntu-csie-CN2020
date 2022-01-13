package transfer

import (
	"log"
	"net"
)

func Recv(conn net.Conn) {
	buf := make([]byte, 1024)
	crrBuf := make([]byte, 0)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			log.Println("error =", err, "connection closing")
			break
		}

		crrBuf = append(crrBuf, buf[:n]...)
		log.Println("Buffer read [", string(crrBuf), "]")
		crrBuf = crrBuf[:0]
	}
}

func RecvMsg(conn net.Conn) string{
	buf := make([]byte, 1024)
	crrBuf := make([]byte, 0)

	n, err := conn.Read(buf)

	if err != nil {
		log.Println("error =", err, "connection closing")
	}

	crrBuf = append(crrBuf, buf[:n]...)
	log.Println("Buffer read [", string(crrBuf), "]")
	return string(crrBuf)
}
