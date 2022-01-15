package transfer

import(
	"net"
	"log"
	"bufio"
	//"time"
)

func Send(conn net.Conn, msg string){
	log.Println("msg =", msg)
	log.Println("12")

	w := bufio.NewWriter(conn)
	log.Println("14")
	n, err := w.WriteString(msg);	
	if err != nil {
		log.Println("write wrong");
		return
	}
	log.Println(n, "bytes sent")
	log.Println("18")
	
	w.Flush()
	log.Println("finish sending [", msg, "]")
}