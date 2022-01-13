package main

import (
	//"./transfer"
	"net"
	"log"
)

func Chat(conn net.Conn, name string){

	log.Println("listing friend in func")

	friendList := ReadFriendList(name)
	SendStringSlice(conn, friendList)

	
}