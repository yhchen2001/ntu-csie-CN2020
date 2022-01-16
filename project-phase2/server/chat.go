package main

import (
	//"./transfer"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	//"fmt"
	"os"
	"time"

	"bytes"

	"./transfer"
)

func Chat(conn net.Conn, myName string) {

	log.Println("listing friend in func")

	friendList := ReadFriendList(myName)
	//time.Sleep(time.Second * 1)
	log.Println("connection = ", conn)
	SendStringSlice(conn, friendList)
	log.Println("chat finish sending friendlist = ", friendList)

	if len(friendList) == 0 {
		log.Println("no friend, return")
		return
	}
	//unknown := transfer.RecvMsg(conn)
	//log.Println("this is =", unknown)

	friendName := transfer.RecvMsg(conn)
	log.Println("friend =", friendName, "is going to chat")
	openChatRoom(conn, myName, friendName)
}

func scanLastNonEmptyLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Set advance to after our last line
	if atEOF {
		advance = len(data)
	} else {
		// data[advance:] now contains a possibly incomplete line
		advance = bytes.LastIndexAny(data, "\n\r") + 1
	}
	data = data[:advance]

	// Remove empty lines (strip EOL chars)
	data = bytes.TrimRight(data, "\n\r")

	// We have no non-empty lines, so advance but do not return a token.
	if len(data) == 0 {
		return advance, nil, nil
	}

	token = data[bytes.LastIndexAny(data, "\n\r")+1:]
	return advance, token, nil
}

func getPath(myName string, friendName string) string {
	var roomName string // roomname 是小＋大
	if myName < friendName {
		roomName = myName + friendName
	} else {
		roomName = friendName + myName
	}

	dir := "./data/chat_history/"
	path := dir + roomName + ".txt"
	return path
}

func appendMsg(path string, line string){
	wrF, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_SYNC|os.O_CREATE, 0644)
	_, _ = wrF.Seek(0, 1)
	if err != nil {
		log.Println("write history file open err :", err)
	}
	if _, err := wrF.WriteString(line); err != nil {
		log.Println("write history file error", err)
	}
	wrF.Close()
}

func monitor(conn net.Conn, myName string, friendName string, exit *bool){
	path := getPath(myName, friendName)
	f, err := os.OpenFile(path, os.O_SYNC|os.O_RDONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil{
		log.Println("monitor open file fail error: ", err)
	}
	oldPos, _ := f.Seek(0, 2)
	log.Println("oldPos =", oldPos)
	for !(*exit){
		newPos, _ := f.Seek(0, 2)
		//log.Println("newPos =", newPos)
		time.Sleep(time.Millisecond * 100)
		if newPos > oldPos{
			log.Printf("new pos = %d , old pos = %d\n", newPos, oldPos)

			f.Seek(oldPos, 0)
			r := bufio.NewReader(f)
			newLine , _ := r.ReadString('\n')
			log.Printf("from %s new line = %s", myName, newLine)
			oldPos = newPos

			s := strings.Split(newLine, ": ")
			log.Println(s)

			if s[0] != myName{
				log.Println("sending to " ,myName, newLine)
				/*w := bufio.NewWriter(conn)
				w.WriteString(newLine)
				w.Flush()*/
				transfer.Send(conn, newLine)
				log.Printf("finish sending %s to %s\n", newLine ,myName)

			}
		}
	}
	log.Println("finish monitoring")
}

func startChat(conn net.Conn, myName string, friendName string) {
	path := getPath(myName, friendName)
	log.Println("pathname =", path)
	r := bufio.NewReaderSize(conn, 1024)

	log.Println("starting chat")

	exit := false
	go monitor(conn, myName, friendName, &exit)

	for {
		line, err := r.ReadString('\n')
		fmt.Printf("read %s", line)
		if err != nil {
			log.Println("error :", err)
			exit = true
			break
		}
		if line == "exit()\n" || line == "exit()" {
			log.Println("exit, return")
			time.Sleep(time.Millisecond * 50)
			transfer.Send(conn, "exit()\n")
			exit = true
			return
		}
		log.Print("shoule write ", line)
		appendMsg(path, myName + ": " +line)
	}
}

func openChatRoom(conn net.Conn, myName string, friend string) {

	log.Println("start writing")
	sendHistory(conn, myName, friend)
	startChat(conn, myName, friend)
}

func sendHistory(conn net.Conn, myName string, friendName string) {
	path := getPath(myName, friendName)
	log.Println("path ==", path)

	if err := os.MkdirAll("./data/chat_history" , os.ModePerm); err != nil {
		log.Println("mk history err: ", err)
	}
	f, err := os.OpenFile(path, os.O_SYNC|os.O_APPEND|os.O_CREATE, 0644)
	defer f.Close()
	log.Println("err file:", err)

	r := bufio.NewReaderSize(f, 1024)
	w := bufio.NewWriterSize(conn, 1024)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				w.WriteString("EOF\n")
				w.Flush()
				log.Println("eof, return")
				break
			}
			log.Println("error :", err)
			time.Sleep(time.Second)
		}
		fmt.Printf("read %s\n", line)
		w.WriteString(line)
		w.Flush()
	}

}
