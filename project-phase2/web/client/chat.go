package client

import (
	"bufio"
	"fmt"
	//"sync"

	//"io"
	//"os"

	"./transfer"

	//"io"
	"log"
	"net"
	"net/http"
	"strconv"
	//"time"
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

var GlobalChatFriendList []string = make([]string, 0)
var GlobalFriendCount int

func WebChatChooseFriend(conn net.Conn, w http.ResponseWriter) {
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
	for i, friend := range GlobalChatFriendList {
		fmt.Fprintf(w, "<input type=\"submit\" value=\"%s\" name=\"%s\"><br>", friend, strconv.Itoa(i))
	}
	//fmt.Fprintln(w, "<input type=\"text\" value=\"rand\" name=\"6\"><br>")
	fmt.Fprintln(w, "</form>")
}

func WebChatChooseFriendPost(conn net.Conn, w http.ResponseWriter, friendCount *int) {

}

var Monitor = false
var GobalMinitorAlive = true
var GlobalMonitorShouldExit = false

func WebStartChat(conn net.Conn, line string) {

	log.Println("starts chating, type exit() to quit")
	GlobalChatHistory = append(GlobalChatHistory, line)
	log.Println(GlobalChatHistory)

	if Monitor == false {
		log.Println("starting minotor!!!!!!!!")
		go webMonitor(conn)
		Monitor = true
	}

	w := bufio.NewWriterSize(conn, 1024)
	w.WriteString(line + "\n")
	w.Flush()
	log.Println("in writer, writed:", line)
	log.Println(line)
	if line == "exit()\n" || line == "exit()" {
		fmt.Println("exit chat room")
		GlobalMonitorShouldExit = true
		log.Println("waiting for monitor to finish")
		for GobalMinitorAlive {
		}
		log.Println("monitor finished detected return")
		return
	}
}

func startChat(conn net.Conn) {
	log.Println("none")
}

func monitor(conn net.Conn, exit *bool, finish *bool) {
	log.Println("start monitoring")
	for !(*exit) {
		msg, err := transfer.RecvMsg(conn)
		if err != nil {
			log.Println("monitor read error, break :", err)
		}
		fmt.Print(msg)
	}
	log.Println("finish monitoring")
	*finish = true
	Monitor = false
}

func webMonitor(conn net.Conn) {
	log.Println("start monitoring")
	for !(GlobalMonitorShouldExit) {
		msg, err := transfer.RecvMsg(conn)
		if err != nil {
			log.Println("monitor read error, break :", err)
		}
		GlobalChatHistory = append(GlobalChatHistory, msg)
		fmt.Println("monitor recieve ", msg)
	}
	log.Println("ending monitoring")
	GobalMinitorAlive = false
}

var GlobalChatHistory []string

func WebDisplayHistory(conn net.Conn, w http.ResponseWriter, tarName string) {

	log.Println("before read")
	transfer.Send(conn, tarName)
	r := bufio.NewReaderSize(conn, 1024)
	//tmpbuf := make([]byte, 0)
	for {
		//log.Println("start")
		line, err := r.ReadString('\n')
		if err != nil {
			log.Println("display err : ", err)
			break
		}

		if line == "EOF\n" {
			log.Print("EOF")
			break
		}
		fmt.Fprintln(w, line, "<br>")
		fmt.Println(line, "<br>")
		GlobalChatHistory = append(GlobalChatHistory, line)
	}
}

func displayHistory(conn net.Conn) {

	log.Println("none")
}
