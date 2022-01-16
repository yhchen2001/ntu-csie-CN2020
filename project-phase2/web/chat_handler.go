package main

import (
	"fmt"
	"net/http"

	//"./client"
	"./client"
	"./client/transfer"
	"log"
	//"time"

	"strconv"
)

var GlobalTarName string
var GetHistory = false
var ChatTarName string

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
		r.ParseForm()
		i := 0
		for ; i < client.GlobalFriendCount; i++ {
			tar := r.FormValue(strconv.Itoa(i))
			log.Println("i = ", tar)
			if len(tar) > 0 {
				ChatTarName = client.GlobalChatFriendList[i]
				log.Println(ChatTarName)
			}
		}
		//test := r.FormValue("6")
		//log.Println("test ", test)
		//fmt.Fprintln(w, "<h2>click home to return<h2>")

		fmt.Fprintf(w,
			"<form action=\"/chatroom\" method=\"GET\">"+
				"<input type=\"submit\" value=\"chatroom\"><br>"+
				"</form>")
		GetHistory = false
		client.GlobalChatHistory = nil
		client.GobalMinitorAlive = true
		client.GlobalMonitorShouldExit = false
		client.Monitor = false
	}
}

func ChatRoomHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprint(w, "<h2>chat room</h2>")
	log.Println("in chat room handler")

	if r.Method == "GET" {

		if GetHistory == false{
			fmt.Println("start getting history")
			client.WebDisplayHistory(Conn, w, ChatTarName)
			GetHistory = true
			fmt.Fprintf(w,
				"<form method=\"GET\">"+
					"type exit() to leave chatroom <input type=\"text\" name=\"line\">"+
					"<input type=\"submit\" value=\"enter\"><br>"+
					"</form>")
		}else{
			r.ParseForm()
			tarLine := r.FormValue("line")


			client.WebStartChat(Conn, tarLine)
			log.Println("in handler, ", client.GlobalChatHistory)

			for _, line := range client.GlobalChatHistory{
				fmt.Fprintln(w, line, "<br>")
			}
			if tarLine == "exit()"{
				BackHome(w)
				return
			}

			fmt.Fprintf(w,
				"<form method=\"GET\">"+
					"type exit() to leave chatroom <input type=\"text\" name=\"line\">"+
					"<input type=\"submit\" value=\"enter\"><br>"+
					"</form>")

		}
	}

	if r.Method == "POST" {

	}
	
}
