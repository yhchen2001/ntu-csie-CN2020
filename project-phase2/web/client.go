package main

import (
	"fmt"
	"log"
	"net/http"
	//"html/template"
)


func backHome(w http.ResponseWriter) {
	fmt.Fprintf(w,
		"<form action=\"/home\" method=\"GET\">"+
			"<input type=\"submit\" value=\"back\">"+
			"</form>")
}

func listFriendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>these are your friends</h1>")

	friendList := []string{"usr1", "usr2", "usr3"}
	for i, friend := range friendList {
		fmt.Fprintf(w, "(%d) %s<br>", i, friend)
	}

	backHome(w)
}

func addFriendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>these are friends you can add</h1>")

	fmt.Fprintf(w,
		"<form method=\"POST\">"+
			"choose a friend to add<input type=\"text\" name=\"target\"><br>"+
			"<input type=\"submit\" value=\"submit\">"+
			"</form>")
	
	backHome(w)
	
	if r.Method == "POST"{
		log.Println("post")
		r.ParseForm()
		target := r.FormValue("target")
		log.Println("target = ", target)
	}
}

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func deleteFriendHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<h1>these are friends</h1>")

	friendList := []string{"usr1", "usr2", "usr3"}
	for i, friend := range friendList {
		fmt.Fprintf(w, "(%d) %s<br>", i, friend)
	}
	fmt.Fprintf(w,
		"<form method=\"POST\">"+
			"choose a friend to delete<input type=\"text\" name=\"target\"><br>"+
			"<input type=\"submit\" value=\"submit\">"+
			"</form>")
	
	backHome(w)
	
	if r.Method == "POST"{
		log.Println("post")
		r.ParseForm()
		target := r.FormValue("target")
		log.Println("target = ", target)
	}

}

func chooseChatFriendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>choose a friend to chat</h1>")

	friendList := []string{"usr1", "usr2", "usr3"}
	for i, friend := range friendList {
		fmt.Fprintf(w, "(%d) %s<br>", i, friend)
		
		fmt.Fprintf(w,
			"<form method=\"POST\">"+
				"%s : <input type=\"submit\" name=\"%s\"><br>"+
				"<input type=\"submit\" value=\"submit\"><br>"+
				"</form>", friend, friend)
	}
	
	backHome(w)
	
	if r.Method == "POST"{
		log.Println("post")
		r.ParseForm()
		target := r.FormValue("target")
		log.Println("target = ", target)

		fmt.Fprintf(w,
			"<form action=\"/chatroom\" method=\"GET\">"+
				"<input type=\"submit\" value=\"start chat\">"+
				"</form>")
	}
}

func chatRoomHandler(w http.ResponseWriter, r * http.Request){
	fmt.Fprintf(w, "<h1>chat room with some user</h1>")

	fmt.Fprintf(w,
		"<form method=\"POST\">"+
			"choose a friend to chat with<input type=\"text\" name=\"target\"><br>"+
			"<input type=\"submit\" value=\"submit\">"+
			"</form>")
		
}

func loginMainHandler(w http.ResponseWriter, r * http.Request){
	
}

func main() {
	fmt.Printf("Starting server at port 8080\n")

	http.HandleFunc("/listfriend", listFriendHandler)
	http.HandleFunc("/addfriend", addFriendHandler)
	http.HandleFunc("/deletefriend", deleteFriendHandler)
	http.HandleFunc("/chooseChatFriend", chooseChatFriendHandler)
	http.HandleFunc("/chatRoom", chatRoomHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
