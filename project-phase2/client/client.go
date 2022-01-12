package main

import (
	"fmt"
	"net"

	"./utils"
)

const(
	server = "127.0.0.1:8080"
)

func send(conn net.Conn, msg string){
	utils.Log("msg =", msg)
	conn.Write([]byte(msg))
	utils.Log("finish sending [", msg, "]")
}

func main(){

	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	utils.CheckErrFatal(err)

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	utils.CheckErrFatal(err)

	for{

		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		
		if err != nil{
			utils.Log("error =", err, "connection closing")
			return
		}

		var msg string
		fmt.Scanln(&msg)

		send(conn, msg)
	}
}