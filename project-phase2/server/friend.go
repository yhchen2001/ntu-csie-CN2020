package main

import (
	//"fmt"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"./transfer"
	//"strconv"
	"strings"
)

func ReadFriendList(name string) []string {
	friendList := make([]string, 0)

	myDir := "./data/" + name
	if err := os.MkdirAll(myDir, os.ModePerm); err != nil {
		log.Println("mkdir err: ", err)
	}

	myFriendFile := myDir + "/allFriends.txt"
	f, err := os.OpenFile(myFriendFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("read myfriend file error", err)
	}
	r := bufio.NewReader(f)
	defer f.Close()

ReadLoop:
	for {
		var crr_name string
		_, err := fmt.Fscanf(r, "%s", &crr_name)
		if err != nil {
			if err == io.EOF {
				break ReadLoop
			}
			continue
		}
		log.Printf("%s\n", crr_name)
		friendList = append(friendList, crr_name)
	}
	return friendList
}

func isFriend(user User, friendList []string) bool {
	for _, friend := range friendList {
		if user.Name == friend {
			return true
		}
	}
	return false
}

func writeFriendList(name string, friendList []string) {

	myDir := "./data/" + name
	myFriendFile := myDir + "/allFriends.txt"

	wrf, err := os.OpenFile(myFriendFile, os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("write myfriend file error", err)
	}
	defer wrf.Close()

	for _, friend := range friendList {
		if _, err := wrf.WriteString(friend + "\n"); err != nil && len(friend) > 0 {
			log.Println("write to friendlist error:", err)
		}
	}
	fmt.Println("finish writing")
}

func stripSpace(list []string) []string {
	str1 := strings.Join(list, " ")
	str2 := strings.Trim(str1, " ")
	return strings.Split(str2, " ")
}

func SendStringSlice(conn net.Conn, sl []string) {
	s := strings.Join(sl, " ")
	if len(sl) == 0 {
		s = "-1"
	}
	transfer.Send(conn, s)
}

func AddFriend(conn net.Conn, myName string) {
	log.Println("adding friend in func")

	myFriendList := ReadFriendList(myName)
	log.Println("friendlist is ", myFriendList)
	userList := ReadUserInfo()
	log.Println("user info length is ", len(userList))
	notFriendList := make([]string, 0)

	for _, user := range userList {
		if isFriend(user, myFriendList) == false && user.Name != myName {
			notFriendList = append(notFriendList, user.Name)
			log.Println(user.Name, "is not my friend")
		}
	}

	if len(notFriendList) == 0 {
		log.Println("no more friends to add")
		return
	}

	SendStringSlice(conn, notFriendList)
	log.Println("all notfriend send finish")

	newFriendName := transfer.RecvMsg(conn)
	log.Println("friend =", newFriendName, "is being added")

	friendsFriendList := ReadFriendList(newFriendName)
	friendsFriendList = append(friendsFriendList, myName)
	writeFriendList(newFriendName, friendsFriendList)

	log.Println("friendlist before is ", myFriendList)
	myFriendList = append(myFriendList, newFriendName)
	log.Println("new friend list =", myFriendList)

	writeFriendList(myName, myFriendList)
}


func DeleteFriend(conn net.Conn, myName string) {

	log.Println("deleting friend in func")
	friendList := ReadFriendList(myName)
	SendStringSlice(conn, friendList)

	if len(friendList) == 0 {
		log.Println("no friends to delete, return")
		return
	}

	notFriendName := transfer.RecvMsg(conn)

	log.Println("old friend list =", friendList)
	if len(friendList) == 1 {
		friendList = nil
	} else {
		for i, friend := range friendList {
			if friend == notFriendName {
				friendList[i] = friendList[len(friendList)-1]
				friendList = friendList[:len(friendList)-1]
			}
		}
	}
	log.Println("new friend list =", friendList)
	writeFriendList(myName, friendList)

	notfriendsFriendList := ReadFriendList(notFriendName)

	log.Println("old not friend list =", notfriendsFriendList)
	if len(notfriendsFriendList) == 1 {
		notfriendsFriendList = nil
	} else {
		for i, friend := range notfriendsFriendList {
			if friend == myName {
				notfriendsFriendList[i] = notfriendsFriendList[len(notfriendsFriendList)-1]
				notfriendsFriendList = notfriendsFriendList[:len(notfriendsFriendList)-1]
			}
		}
	}
	log.Println("new not friend list =", friendList)
	writeFriendList(notFriendName, notfriendsFriendList)
}

func ListFriend(conn net.Conn, name string) {
	log.Println("listing friend in func")

	friendList := ReadFriendList(name)
	SendStringSlice(conn, friendList)
}
