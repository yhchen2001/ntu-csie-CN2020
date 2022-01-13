package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"./transfer"
	"sync"
)





var WriteInfoMutex sync.Mutex
var ReadInfoMutex sync.Mutex
var OnlineUsers = make([]User, 100)

type User struct{
	Name string
	Pass string
}

func RemoveUser(slice []User, i int) []User {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func readUserInfo() []User{
	ReadInfoMutex.Lock()
	defer ReadInfoMutex.Unlock()

	userList := make([]User, 100)
	f, err := os.Open("./data/user_info.txt")
	if err != nil {
		log.Println("open user_info fail")
	}
	r := bufio.NewReader(f)

ReadLoop:
	for{
		var crr_name string
		var crr_pass string
		_, err := fmt.Fscanf(r, "%s%s", &crr_name, &crr_pass)
		if err != nil{
			if err == io.EOF{
				break ReadLoop
			}
			log.Println("read user_info error :", err)
		}
		log.Printf("%s %s\n", crr_name, crr_pass)
		userList = append(userList, User{crr_name, crr_pass})
	}
	f.Close()
	return userList
}

func checkSignIn(name string, pass string) string{
	WriteInfoMutex.Lock()
	defer WriteInfoMutex.Unlock()

	userList := readUserInfo()

	for _, user := range OnlineUsers{
		if user.Name == name{
			log.Println("user online")
			return "fail"
		}
	}

	log.Println("name pass = ", name, pass)

	for _, user := range userList{
		log.Println("name pass = ", user.Name, user.Pass)
		if user.Name == name && user.Pass == pass{
			OnlineUsers = append(OnlineUsers, User{name, pass})
			return "ok"
		}
	}
	log.Println("username or password incorrect")
	return "fail"
}

func checkSignUp(name string, pass string) string{
	WriteInfoMutex.Lock()
	defer WriteInfoMutex.Unlock()

	userList := readUserInfo()

	for _, user := range userList{
		if user.Name == name{
			return "fail"
		}
	}

	log.Println("start writing")
	wrF, _ := os.OpenFile("./data/user_info.txt", os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0644)
	if _, err := wrF.WriteString(name + " " + pass + "\n"); err != nil {
		log.Println("write error", err)
	}
	log.Println("finish checking")

	OnlineUsers = append(OnlineUsers, User{name, pass})
	return "ok"
}


func SignIn(conn net.Conn) (string, string, string){
	log.Println("start signing in")
	name := transfer.RecvMsg(conn)
	log.Println(name)
	password := transfer.RecvMsg(conn)
	log.Println(password)

	okMsg := checkSignIn(name, password)

	log.Println("start sending ", okMsg)
	transfer.Send(conn, okMsg)
	log.Println("finish sending ", okMsg)
	return name, password, okMsg
}

func SignUp(conn net.Conn) (string, string, string) {
	log.Println("start signing up")
	name := transfer.RecvMsg(conn)
	password := transfer.RecvMsg(conn)

	okMsg := checkSignUp(name, password)

	log.Println("start sending ", okMsg)
	transfer.Send(conn, okMsg)
	log.Println("finish sending ", okMsg)
	return name, password, okMsg
}
