package main

import (
	"fmt"
	"log"
	"net"

	"./transfer"
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

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn){

	utils.Log("Handling user...")
	defer conn.Close()
	var name string
SignLoop:
	for{
		msg := transfer.RecvMsg(conn)
		utils.Log("Buffer read [", msg, "]")
		switch msg{
			case "1" :	
				if tmpname, _, okMsg := SignIn(conn); okMsg == "ok"{
					name = tmpname
					break SignLoop
				}
				log.Println("sign in fail")
			case "2" :
				if tmpname, _, okMsg := SignUp(conn); okMsg == "ok"{
					name = tmpname
					break SignLoop
				}
				log.Println("sign up fail")
			default :
				fmt.Println("input format wrong")
		}
	}

	fmt.Println("finish logging in, name = ", name)

	defer func(){
		for i, user := range OnlineUsers{
			if user.Name == name{
				RemoveUser(OnlineUsers, i)
				break
			}
		}
	}()

	for{
		action := transfer.RecvMsg(conn)

		switch action{
			case "1":
				log.Println("Listing friend")
				ListFriend(conn)
			case "2":
				log.Println("Adding a friend")
				AddFriend(conn)
			case "3":
				log.Println("Deleting a friend")
				DeleteFriend(conn)
			case "4":
				log.Println("chat")
				Chat(conn)
			default:
				log.Println("wrong input format")
		}
	}

}

