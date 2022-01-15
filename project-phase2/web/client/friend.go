package client

import (
	//"fmt"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"strconv"

	"./transfer"
)

func ReceiveStringSlice(conn net.Conn) []string{
	msg, _ := transfer.RecvMsg(conn)

	if msg == "-1" {
		return make([]string, 0)
	}

	sl := strings.Split(msg, " ")
	return sl
}

func ScanAndCheckChoice(list []string) (int){
	var tarName string
	var tarNum int
	for{
		fmt.Scanf("%s", &tarName)
		i, err := strconv.Atoi(tarName)
		tarNum = i
		if err != nil {
			fmt.Println("Input format wrong")
			continue
		}
		if tarNum >= len(list) || tarNum < 0{
			fmt.Printf("Number not available")
			continue
		}
		break
	}
	return tarNum
}

var globalNotFriendList []string

func WebAddFriend(conn net.Conn, w http.ResponseWriter){
	log.Println("adding friend in func")

	fmt.Fprintf(w, "<h1>these are friends you can add</h1>")

	globalNotFriendList = ReceiveStringSlice(conn)
	if len(globalNotFriendList) == 0{
		fmt.Fprintln(w, "you don't have friend to delete")
		return
	}
	for i, notFriend := range globalNotFriendList{
		fmt.Fprintf(w, "(%d) %s<br>", i, notFriend)
	}	
}

func WebAddFriendSecond(conn net.Conn, w http.ResponseWriter, tarName string){
	tarNum, _ := strconv.Atoi(tarName)
	fmt.Fprintln(w, "<h2>you have added = " ,globalNotFriendList[tarNum], "<h2>")
	transfer.Send(conn, globalNotFriendList[tarNum])
}

func AddFriend(conn net.Conn){
	log.Println("adding friend in func")

	notFriendList := ReceiveStringSlice(conn)
	if len(notFriendList) == 0{
		fmt.Println("you don't have friend to delete")
		return
	}
	
	fmt.Println("those are friends you can add")
	for i, notFriend := range notFriendList{
		fmt.Printf("(%d) %s\n", i, notFriend)
	}
	fmt.Printf("choose a friend to add :")
	
	tarNum := ScanAndCheckChoice(notFriendList)
	
	fmt.Println("target = " ,notFriendList[tarNum])
	transfer.Send(conn, notFriendList[tarNum])	
}

var globalFriendList []string

func WebDeleteFriendGet(conn net.Conn, w http.ResponseWriter){
	log.Println("deleting friend in func")
	globalFriendList = ReceiveStringSlice(conn)

	if len(globalFriendList) == 0{
		fmt.Fprintln(w, "you don't have friend to delete<br>")
		return
	}

	fmt.Fprintln(w, "those are your friend<br>")
	for i, notFriend := range globalFriendList{
		fmt.Fprintf(w, "(%d) %s<br>", i, notFriend)
	}
}

func WebDeleteFriendPost(conn net.Conn, w http.ResponseWriter, tarName string){	
	
	tarNum, _ := strconv.Atoi(tarName)
	fmt.Fprintf(w, "you have deleted target %s<br>" ,globalFriendList[tarNum])
	transfer.Send(conn, globalFriendList[tarNum])
}

func DeleteFriend(conn net.Conn){
	log.Println("deleting friend in func")
	friendList := ReceiveStringSlice(conn)

	if len(friendList) == 0{
		fmt.Println("you don't have friend to delete")
		return
	}

	fmt.Println("those are your friend")
	for i, notFriend := range friendList{
		fmt.Printf("(%d) %s\n", i, notFriend)
	}
	fmt.Printf("choose a friend to delete :")	
	
	tarNum := ScanAndCheckChoice(friendList)
	fmt.Println("target = " ,friendList[tarNum])
	transfer.Send(conn, friendList[tarNum])
}

func WebListFriend(conn net.Conn, w http.ResponseWriter){
	log.Println("listing friend in func")

	friendList := ReceiveStringSlice(conn)

	for i, friend := range(friendList){
		fmt.Fprintf(w, "(%d) %s<br>", i, friend)
	}
	fmt.Println("<br>")

}

func ListFriend(conn net.Conn){
	log.Println("listing friend in func")


	friendList := ReceiveStringSlice(conn)

	fmt.Println("your friends are")
	for i, friend := range(friendList){
		fmt.Printf("(%d) %s<br>", i, friend)
	}
	fmt.Println("<br>")

}