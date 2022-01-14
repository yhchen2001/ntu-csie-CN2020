package main

import (
	//"fmt"
	"fmt"
	"log"
	"net"
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

func ListFriend(conn net.Conn){
	log.Println("listing friend in func")


	friendList := ReceiveStringSlice(conn)

	fmt.Println("your friends are")
	for i, friend := range(friendList){
		fmt.Printf("(%d) %s\n", i, friend)
	}
	fmt.Println("")

}