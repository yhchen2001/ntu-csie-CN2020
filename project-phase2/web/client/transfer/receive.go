package transfer

import (
	"io"
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

func RecvMsg(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	crrBuf := make([]byte, 0)
	log.Println("30")
	n, err := conn.Read(buf)
	log.Println("32")
	if err != nil {
		log.Println("error =", err, "connection closing")
		return "", io.EOF
	}
	log.Println("37")

	crrBuf = append(crrBuf, buf[:n]...)
	log.Println("Buffer read [", string(crrBuf), "]")
	return string(crrBuf), nil
}
