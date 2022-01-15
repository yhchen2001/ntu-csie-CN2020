package main

import (
	//"bufio"
	"fmt"
	"log"
	"net/http"

	"./client"
	"net"

	//"./client/transfer"
	"./client/transfer"
	"./client/utils"
	//"html/template"
)

const (
	server   = "127.0.0.1:8080"
	homeMsg  = "(1) list friends\n(2) Add a friend\n(3) Delete a friend\n(4) Chat\nyour action: "
	loginMsg = "(1) sign in\n(2) sign up\nyour action: "
)

func backHome(w http.ResponseWriter) {
	fmt.Fprintf(w,
		"<form action=\"/home\" method=\"GET\">"+
			"<input type=\"submit\" value=\"back\"><br>"+
			"</form>")
}

func listFriendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html")
		transfer.Send(Conn, "1")
		fmt.Fprintf(w, "<h1>these are your friends</h1>")
		client.WebListFriend(Conn, w)
		backHome(w)
	}
}

func addFriendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.Method == "GET" {
		transfer.Send(Conn, "2")
		client.WebAddFriend(Conn, w)
		fmt.Fprintf(w,
			"<form method=\"POST\">"+
				"choose a friend to add<input type=\"text\" name=\"target\"><br>"+
				"<input type=\"submit\" value=\"submit\">"+
				"</form>")
	}

	if r.Method == "POST" {
		log.Println("post")
		r.ParseForm()
		target := r.FormValue("target")
		log.Println("target = ", target)
		client.WebAddFriendSecond(Conn, w, target)
		//fmt.Fprintln(w, "<h2>click home to return<h2>")
		backHome(w)
	}
}

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func deleteFriendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method == "GET" {
		transfer.Send(Conn, "3")
		log.Println("deleting friend in get func")

		client.WebDeleteFriendGet(Conn, w)
		fmt.Fprintf(w,
			"<form method=\"POST\">"+
				"choose a friend to delete<input type=\"text\" name=\"target\"><br>"+
				"<input type=\"submit\" value=\"submit\">"+
				"</form>")
	}

	if r.Method == "POST" {
		log.Println("post")
		r.ParseForm()
		target := r.FormValue("target")
		client.WebDeleteFriendPost(Conn, w, target)
		backHome(w)
	}

}

var Name string

var tcpAddr, _ = net.ResolveTCPAddr("tcp4", server)
var Conn, _ = net.DialTCP("tcp4", nil, tcpAddr)

func main() {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	utils.CheckErrFatal(err)

	Conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	utils.CheckErrFatal(err)
	defer Conn.Close()

	log.Println("111  !!!!!!!", Conn)

	fmt.Printf("Starting server at port 8080\n")

	http.HandleFunc("/listfriend", listFriendHandler)
	http.HandleFunc("/addfriend", addFriendHandler)
	http.HandleFunc("/deletefriend", deleteFriendHandler)
	http.HandleFunc("/signin", SignInHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	if err := http.ListenAndServe("", nil); err != nil {
		log.Fatal(err)
	}
}
