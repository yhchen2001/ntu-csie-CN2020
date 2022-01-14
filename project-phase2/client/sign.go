package main

import (
	"fmt"
	"log"
	"net"

	"./transfer"
	//"bufio"
	//"./utils"
)

func SignIn(conn net.Conn) string{
	name := ""
	password := ""

	fmt.Println("enter your name (<30 word): ")
	fmt.Scanln(&name)
	transfer.Send(conn, name)

	fmt.Println("enter your password (<30 word): ")
	fmt.Scanln(&password)
	transfer.Send(conn, password)



	fmt.Println("start recieving ok msg ")
	msg, _ := transfer.RecvMsg(conn)
	fmt.Println("end recieving ok msg ")

	switch msg {
	case "ok":
		fmt.Println("sign in success")
		return "ok"
	case "fail":
		fmt.Println("username or password incorrect, please try again")
	default:
		log.Printf("server response [%s] is wrong", msg)
	}

	return "fail"
}

func SignUp(conn net.Conn) string {

	name := ""
	password := ""

	fmt.Println("enter your name (<30 word): ")
	fmt.Scanln(&name)
	transfer.Send(conn, name)

	fmt.Println("enter your password (<30 word): ")
	fmt.Scanln(&password)
	transfer.Send(conn, password)


	msg, _ := transfer.RecvMsg(conn)

	switch msg {
	case "ok":
		fmt.Println("sign in success")
		return "ok"
	case "fail":
		fmt.Println("username or password incorrect, please try again")
	default:
		log.Printf("server response [%s] is wrong", msg)
	}

	return "fail"
}
