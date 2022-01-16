package main

import (
	//"bufio"
	"fmt"
	"log"
	"net/http"

	//"./client"
	"net"

	//"./client/transfer"
	//"./client/transfer"
	//"./client/utils"
	//"html/template"
)

const (
	server   = "127.0.0.1:8080"
	homeMsg  = "(1) list friends\n(2) Add a friend\n(3) Delete a friend\n(4) Chat\nyour action: "
	loginMsg = "(1) sign in\n(2) sign up\nyour action: "
)

func BackHome(w http.ResponseWriter) {
	fmt.Fprintf(w,
		"<form action=\"/home\" method=\"GET\">"+
			"<input type=\"submit\" value=\"back\"><br>"+
			"</form>")
}

var Name string

var tcpAddr, _ = net.ResolveTCPAddr("tcp4", server)
var Conn, _ = net.DialTCP("tcp4", nil, tcpAddr)

func main() {
	defer Conn.Close()

	log.Println("111  !!!!!!!", Conn)

	fmt.Printf("Starting server at port 8080\n")

	http.HandleFunc("/listfriend", listFriendHandler)
	http.HandleFunc("/addfriend", addFriendHandler)
	http.HandleFunc("/deletefriend", deleteFriendHandler)
	http.HandleFunc("/signin", SignInHandler)
	http.HandleFunc("/signup", SignUpHandler)
	http.HandleFunc("/chooseChatFriend", ChooseChatFriendHandler)
	http.HandleFunc("/chatroom", ChatRoomHandler)
	http.HandleFunc("/profile", ProfileHandler)
	http.HandleFunc("/changePassword", ChangePasswordHandler)
	

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	if err := http.ListenAndServe("", nil); err != nil {
		log.Fatal(err)
	}
}
