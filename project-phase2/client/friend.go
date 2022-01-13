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

func receiveStringSlice(conn net.Conn) []string{
	msg := transfer.RecvMsg(conn)

	if msg == "-1" {
		return make([]string, 0)
	}

	sl := strings.Split(msg, " ")
	return sl
}

func AddFriend(conn net.Conn){
	log.Println("adding friend in func")

	notFriendList := receiveStringSlice(conn)
	if len(notFriendList) == 0{
		fmt.Println("you don't have friend to delete")
		return
	}
	
	fmt.Println("those are friends you can add")
	for i, notFriend := range notFriendList{
		fmt.Printf("(%d) %s\n", i, notFriend)
	}
	fmt.Printf("choose a friend to add :")
	
	var tarNotFriend string
	var tarNum int
	for{
		fmt.Scanf("%s", &tarNotFriend)
		i, err := strconv.Atoi(tarNotFriend)
		tarNum = i
		if err != nil {
			fmt.Println("Input format wrong")
			continue
		}
		break
	}
	fmt.Println("tarnotfriend = " ,tarNotFriend)
	transfer.Send(conn, notFriendList[tarNum])	
}

func DeleteFriend(conn net.Conn){
	log.Println("deleting friend in func")
	friendList := receiveStringSlice(conn)

	if len(friendList) == 0{
		fmt.Println("you don't have friend to delete")
		return
	}

	fmt.Println("those are your friend")
	for i, notFriend := range friendList{
		fmt.Printf("(%d) %s\n", i, notFriend)
	}
	fmt.Printf("choose a friend to delete :")	
	
	var tarNotFriend string
	var tarNum int
	for{
		fmt.Scanf("%s", &tarNotFriend)
		i, err := strconv.Atoi(tarNotFriend)
		tarNum = i
		if err != nil {
			fmt.Println("Input format wrong")
			continue
		}
		break
	}
	fmt.Println("tarnotfriend = " ,tarNotFriend)
	transfer.Send(conn, friendList[tarNum])
}

func ListFriend(conn net.Conn){
	log.Println("listing friend in func")


	friendList := receiveStringSlice(conn)

	fmt.Println("your friends are")
	for i, friend := range(friendList){
		fmt.Printf("(%d) %s\n", i, friend)
	}
	fmt.Println("")

}