package transfer

import(
	"net"
	"log"
	"bufio"
	"time"
)

func Send(conn net.Conn, msg string){
	log.Println("msg =", msg)

	w := bufio.NewWriter(conn)
	if _, err := w.WriteString(msg); err != nil {
		log.Println("write wrong");
		return
	}
	w.Flush()
	time.Sleep(time.Millisecond)
	log.Println("finish sending [", msg, "]")
}