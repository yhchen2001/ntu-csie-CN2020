package main

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

func startChat(conn net.Conn){
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

func displayHistory(conn net.Conn) {

	log.Println("before read")
	r := bufio.NewReaderSize(conn, 1024)

	for{
		//log.Println("start")
		line, err := r.ReadString('\n')
		fmt.Print(line)

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
