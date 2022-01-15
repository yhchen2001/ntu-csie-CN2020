package main

import (
	"fmt"
	"net/http"

	//"./client"
	"./client"
	"./client/transfer"
	"log"

	"strconv"
)

func ChooseChatFriendHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	//var tarNum *int

	if r.Method == "GET" {
		transfer.Send(Conn, "4")
		fmt.Fprintln(w, "<h1>what the hell</h1>")
		client.WebChatChooseFriend(Conn, w)
		log.Println("friendcount =", client.GlobalFriendCount)
		log.Println(client.GlobalChatFriendList)
	}

	if r.Method == "POST" {
		log.Println("post")
		var tarName string
		r.ParseForm()
		i := 0
		for ; i < client.GlobalFriendCount; i++ {
			tar := r.FormValue(strconv.Itoa(i))
			log.Println("i = ", tar)
			if len(tar) > 0 {
				tarName = client.GlobalChatFriendList[i]
				log.Println(tarName)
			}
		}
		//test := r.FormValue("6")
		//log.Println("test ", test)
		transfer.Send(Conn, tarName)
		//fmt.Fprintln(w, "<h2>click home to return<h2>")

		fmt.Fprintf(w,
			"<form action=\"/chatroom\" method=\"GET\">"+
				"<input type=\"submit\" value=\"chatroom\"><br>"+
				"</form>")
	}
}

func ChatRoomHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprint(w, "<h2>chat room<h2>")
}
