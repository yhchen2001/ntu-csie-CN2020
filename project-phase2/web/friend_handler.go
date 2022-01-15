package main

import (
	//"bufio"
	"fmt"
	"log"
	"net/http"

	"./client"
	//"net"

	//"./client/transfer"
	"./client/transfer"
	//"./client/utils"
	//"html/template"
)

func listFriendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html")
		transfer.Send(Conn, "1")
		fmt.Fprintf(w, "<h1>these are your friends</h1>")
		client.WebListFriend(Conn, w)
		BackHome(w)
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
		BackHome(w)
	}
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
		BackHome(w)
	}

}