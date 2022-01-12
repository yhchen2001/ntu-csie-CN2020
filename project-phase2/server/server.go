package main

import (
	"net"
	"./utils"
	_ "github.com/go-sql-driver/mysql"
)


const(
	port = "8080"
)

func main(){
	startServer()
}

func startServer(){
	// open a connection
	mainSock, err := net.Listen("tcp", "localhost:" + port)
	utils.CheckErrFatal(err)
	utils.Log("starting at", ("localhost:" + port))
	defer mainSock.Close()


	utils.Log("waiting for client...")

	for{
		conn, err := mainSock.Accept()
		if err != nil{
			continue
		}
		utils.Log(conn.RemoteAddr().String(), "TCP connection success")

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn){

	utils.Log("Handling connection...")

	buf := make([]byte, 1024)
	crrBuf := make([]byte, 0)

	for{
		n, err := conn.Read(buf)
		
		if err != nil{
			utils.Log("error =", err, "connection closing")
			return
		}

		crrBuf = append(crrBuf, buf[:n]...)
		utils.Log("Buffer read [", string(crrBuf), "]")
		crrBuf = crrBuf[:0]
	}

}

