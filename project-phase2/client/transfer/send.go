package transfer

import(
	"net"
	"log"
	"bufio"
)

func Send(conn net.Conn, msg string){
	log.Println("msg =", msg)

	w := bufio.NewWriter(conn)
	if _, err := w.WriteString(msg); err != nil {
		log.Println("flush wrong");
		return
	}
	w.Flush()
	
	log.Println("finish sending [", msg, "]")
}