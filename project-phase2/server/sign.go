package main

import (
	"./transfer"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

var writeInfoMutex sync.Mutex
var readInfoMutex sync.Mutex
var OnlineUsers = make([]User, 0)

type User struct {
	Name string
	Pass string
	Conn net.Conn
}

func RemoveUser(slice []User, i int) []User {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func ReadUserInfo() []User {
	readInfoMutex.Lock()
	defer readInfoMutex.Unlock()

	userList := make([]User, 0)
	f, err := os.Open("./data/user_info.txt")
	defer f.Close()
	if err != nil {
		log.Println("open user_info fail")
	}
	r := bufio.NewReader(f)

ReadLoop:
	for {
		var crr_name string
		var crr_pass string
		_, err := fmt.Fscanf(r, "%s%s", &crr_name, &crr_pass)
		if err != nil {
			if err == io.EOF {
				break ReadLoop
			}
			continue
		}
		log.Printf("%s %s\n", crr_name, crr_pass)
		userList = append(userList, User{crr_name, crr_pass, nil})
	}
	f.Close()
	return userList
}

func checkSignIn(conn net.Conn, name string, pass string) string {
	writeInfoMutex.Lock()
	defer writeInfoMutex.Unlock()

	userList := ReadUserInfo()

	for _, user := range OnlineUsers {
		if user.Name == name {
			log.Println("user online")
			return "fail"
		}
	}

	log.Println("name pass = ", name, pass)

	for _, user := range userList {
		log.Println("name pass = ", user.Name, user.Pass)
		if user.Name == name && user.Pass == pass {
			OnlineUsers = append(OnlineUsers, User{name, pass, conn})
			return "ok"
		}
	}
	log.Println("username or password incorrect")
	return "fail"
}

func checkSignUp(conn net.Conn, name string, pass string) string {
	writeInfoMutex.Lock()
	defer writeInfoMutex.Unlock()

	userList := ReadUserInfo()

	for _, user := range userList {
		if user.Name == name {
			return "fail"
		}
	}

	log.Println("start writing")
	wrF, _ := os.OpenFile("./data/user_info.txt", os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0644)
	defer wrF.Close()
	if _, err := wrF.WriteString(name + " " + pass + "\n"); err != nil {
		log.Println("write error", err)
	}
	log.Println("finish checking")

	OnlineUsers = append(OnlineUsers, User{name, pass, conn})
	return "ok"
}

func SignIn(conn net.Conn) (string, string, string) {
	log.Println("start signing in")
	name := transfer.RecvMsg(conn)
	log.Println(name)
	password := transfer.RecvMsg(conn)
	log.Println(password)

	okMsg := checkSignIn(conn, name, password)

	log.Println("start sending ", okMsg)
	transfer.Send(conn, okMsg)
	log.Println("finish sending ", okMsg)
	return name, password, okMsg
}

func SignUp(conn net.Conn) (string, string, string) {
	log.Println("start signing up")
	name := transfer.RecvMsg(conn)
	password := transfer.RecvMsg(conn)

	okMsg := checkSignUp(conn, name, password)

	log.Println("start sending ", okMsg)
	transfer.Send(conn, okMsg)
	log.Println("finish sending ", okMsg)
	return name, password, okMsg
}

func ChangePassword(conn net.Conn, name string) {
	log.Println("start changing password")

	newPass := transfer.RecvMsg(conn)
	log.Println("finish receiveing new pass = ", newPass)

	userList := ReadUserInfo()

	readInfoMutex.Lock()
	defer readInfoMutex.Unlock()

	wrf, err := os.OpenFile("./data/user_info.txt", os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0644)
	defer wrf.Close()

	if err != nil {
		log.Println("open user_info fail")
	}

	for _, user := range userList {
		var newLine string
		if user.Name == name {
			newLine = user.Name + " " + newPass + "\n"
		} else {
			newLine = user.Name + " " + user.Pass + "\n"
		}
		log.Print("new line =", newLine)
		wrf.WriteString(newLine)
	}

	log.Println("finish modifying user_info.txt")
}
