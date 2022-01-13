package main

import(
	"net"
	"fmt"
)

func Chat(conn net.Conn){
	friendList := receiveStringSlice(conn)

	fmt.Println("choose a friends to chat with")
	for i, friend := range(friendList){
		fmt.Printf("(%d) %s\n", i, friend)
	}
	fmt.Println("")

}