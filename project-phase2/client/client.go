package main

import (
	"fmt"
	"net"
	//"net/http"

	"./utils"
	"./transfer"
	//"./actions"

	//".web"
)

const(
	server = "127.0.0.1:8080"
	homeMsg = "(1) list friends\n(2) Add a friend\n(3) Delete a friend\n(4) Chat\nyour action: "
	loginMsg = "(1) sign in\n(2) sign up\nyour action: "
)

func main(){

	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	utils.CheckErrFatal(err)

	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	utils.CheckErrFatal(err)
	defer conn.Close()

SigninLoop:	
	for{
		fmt.Printf(loginMsg)
		var msg string
		fmt.Scanln(&msg)
		transfer.Send(conn, msg)

		switch msg{
			case "1" :	
				if SignIn(conn) == "ok"{
					fmt.Println("ok msg =", msg)
					break SigninLoop
				}
			case "2" :
				if SignUp(conn) == "ok"{
					break SigninLoop
				}
			default :
				fmt.Println("input format wrong")
		}
	}
	fmt.Println("finish logging in")

	for{
		fmt.Printf(homeMsg)
		var msg string
		fmt.Scanln(&msg)

		switch msg{
			case "1" :
				ListFriend()
			case "2" :
				AddFriend()
			case "3" :
				DeleteFriend()
			case "4" :
				Chat()
			default :
				utils.Log("wrong format choice")
		}
		
		transfer.Send(conn, msg)
	}
}