package client

import (
	"bufio"
	"fmt"
	//"sync"

	//"io"
	"os"

	"./transfer"

	//"io"
	"log"
	"net"
	"net/http"
	"strconv"
	//"os"
)

func Chat(conn net.Conn) {
	friendList := ReceiveStringSlice(conn)
	if len(friendList) == 0 {
		fmt.Println("you have no friends to chat with")
		return
	}

	fmt.Println("choose a friends to chat with")
	for i, friend := range friendList {
		fmt.Printf("(%d) %s\n", i, friend)
	}
	fmt.Printf("your choice: ")

	tarNum := ScanAndCheckChoice(friendList)

	fmt.Println("target = ", friendList[tarNum])
	transfer.Send(conn, friendList[tarNum])

	// starts chatting
	displayHistory(conn)

	startChat(conn)
}

func WebChatChooseFriend_backup(conn net.Conn, w http.ResponseWriter, friendCount *int, tarName * string) {
	log.Println("start recieving friend")
	
	friendList := ReceiveStringSlice(conn)
	log.Println("end recieving friend")
	if len(friendList) == 0 {
		fmt.Fprintln(w, "you have no friends to chat with<br>")
		return
	}

	fmt.Fprintln(w, "<h3>choose a friends to chat with</h3><br>")
	fmt.Fprintln(w, "<form method=\"POST\">")
	*friendCount = len(friendList)
	for i, friend := range friendList {
		fmt.Fprintf(w, "<input type=\"submit\" value=\"%s\" name=\"%d\"><br>", friend, i)
	}
	fmt.Fprintln(w, "</form>")
}

var GlobalChatFriendList []string = make([]string, 0)
var GlobalFriendCount int

func WebChatChooseFriend(conn net.Conn, w http.ResponseWriter){
	log.Println("choose friend in func")

	fmt.Fprintf(w, "<h2>choose a friend to chat</h2>")
	GlobalChatFriendList = ReceiveStringSlice(conn)
	GlobalFriendCount = len(GlobalChatFriendList)
/*
	if len(GlobalChatFriendList) == 0{
		fmt.Fprintln(w, "you don't have friend to chat")
		return
	}
*/
	fmt.Fprintln(w, "<form method=\"POST\">")
	for i, friend := range GlobalChatFriendList{	
		fmt.Fprintf(w, "<input type=\"submit\" value=\"%s\" name=\"%s\"><br>", friend, strconv.Itoa(i))
	}
	//fmt.Fprintln(w, "<input type=\"text\" value=\"rand\" name=\"6\"><br>")
	fmt.Fprintln(w, "</form>")	
}

func WebChatChooseFriendPost(conn net.Conn, w http.ResponseWriter, friendCount *int){

}

func WebStartChat(conn net.Conn){
	log.Println("starts chating, type exit() to quit")
	w := bufio.NewWriterSize(conn, 1024)
	r := bufio.NewReader(os.Stdin)


	exit := false
	finish := false
	go monitor(conn, &exit, &finish)

	defer func(){
		log.Println("waiting for monitor to finish")
		for !finish{
		}
	}()

	for{
		log.Println("before read")
		line, err := r.ReadString('\n')
		log.Print("after read line =", line)

		w.WriteString(line)
		w.Flush()
		if line == "exit()\n"{
			fmt.Println("exit chat room")
			exit = true
			return
		}
		if err != nil{
			log.Println("error = ", err, " return")
			return
		}
	}
}

func startChat(conn net.Conn){
	log.Println("none")
}

func monitor(conn net.Conn, exit * bool, finish * bool){
	log.Println("start monitoring")
	for !(*exit) {
		msg,err := transfer.RecvMsg(conn)
		if err != nil{
			log.Println("monitor read error, break :", err)
		}
		fmt.Print(msg)
	}
	log.Println("finish monitoring")
	*finish = true
}

func WebDisplayHistory(conn net.Conn, w http.ResponseWriter) {

	log.Println("before read")
	r := bufio.NewReaderSize(conn, 1024)

	for{
		//log.Println("start")
		line, err := r.ReadString('\n')
		fmt.Fprintln(w, line, "<br>")

		if err != nil{
			log.Println("display err : ", err)
			break;
		}
		if line == "EOF\n"{
			log.Print("EOF")
			break;
		}
	}
}

func displayHistory(conn net.Conn) {

	log.Println("none")
}
