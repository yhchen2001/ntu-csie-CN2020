package main

import (
	//"bufio"
	"fmt"
	"log"
	"net/http"
	"time"

	"./client/transfer"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if r.Method == "GET" {

		log.Println("signing in, send \"1\"")
		log.Println("!!!!!!!!!!!!!!!!!!!!", Conn)
		transfer.Send(Conn, "1")

		fmt.Fprintf(w,
			"<form method=\"POST\">"+
				"usrname<input type=\"text\" name=\"name\"><br>"+
				"password<input type=\"text\" name=\"password\"><br>"+
				"<input type=\"submit\" value=\"submit\">"+
				"</form>")
		fmt.Println("end of get")
	}

	if r.Method == "POST" {
		log.Println("post")
		r.ParseForm()
		Name = r.FormValue("name")
		password := r.FormValue("password")
		fmt.Fprintf(w, "name, pass = %s, %s<br>", Name, password)
		log.Println("name, pass =", Name, password, "<br>")

		log.Println("!!!!!!!!!!!!!!!!!!!!", Conn)
		transfer.Send(Conn, Name)
		time.Sleep(time.Millisecond * 10)
		transfer.Send(Conn, password)

		fmt.Println("start recieving ok msg ")
		msg, _ := transfer.RecvMsg(Conn)
		fmt.Println("end recieving ok msg ")

		switch string(msg) {
		case "ok":
			fmt.Fprintln(w, "sign in success")
			BackHome(w)
		case "fail":
			fmt.Fprintln(w, "username or password incorrect, please try again")

			fmt.Fprintf(w, "<form action=\"signin\" method=\"GET\">"+
				"<input type=\"submit\" value=\"try again\">"+
				"</form>")

		default:
			log.Printf("server response [%s] is wrong", msg)
		}
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<h2>signing up</h2>")

	if r.Method == "GET" {
		log.Println("signing up, send \"2\"")
		log.Println("!!!!!!!!!!!!!!!!!!!!", Conn)
		transfer.Send(Conn, "2")

		fmt.Fprintf(w,
			"<form method=\"POST\">"+
				"usrname<input type=\"text\" name=\"name\"><br>"+
				"password<input type=\"text\" name=\"password\"><br>"+
				"<input type=\"submit\" value=\"submit\">"+
				"</form>")
		fmt.Println("end of get")
	}

	if r.Method == "POST" {
		log.Println("post")
		r.ParseForm()
		Name = r.FormValue("name")
		password := r.FormValue("password")
		fmt.Fprintf(w, "name, pass = %s, %s<br>", Name, password)
		log.Println("name, pass =", Name, password, "<br>")

		log.Println("!!!!!!!!!!!!!!!!!!!!", Conn)
		transfer.Send(Conn, Name)
		time.Sleep(time.Millisecond * 10)
		transfer.Send(Conn, password)

		fmt.Println("start recieving ok msg ")
		msg, _ := transfer.RecvMsg(Conn)
		fmt.Println("end recieving ok msg ")

		switch string(msg) {
		case "ok":
			fmt.Fprintln(w, "sign up success")
			BackHome(w)
		case "fail":
			fmt.Fprintln(w, "username or password incorrect, please try again")

			fmt.Fprintf(w, "<form action=\"signup\" method=\"GET\">"+
				"<input type=\"submit\" value=\"try again\">"+
				"</form>")

		default:
			log.Printf("server response [%s] is wrong", msg)
		}
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, "<h3>This is your profile</h3>")

	fmt.Fprintf(w, "<form action=\"changePassword\" method=\"GET\">"+
		"<input type=\"submit\" value=\"Change password\">"+
		"</form>")

}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.Method == "GET" {
		transfer.Send(Conn, "5")
		fmt.Fprintln(w, "<form method=\"POST\">"+
			"enter your new password<input type=\"text\" name=\"newPass\"><br>"+
			"<input type=\"submit\" value=\"submit\">"+
			"</form>")
	}
	if r.Method == "POST" {
		r.ParseForm()
		newPass := r.FormValue("newPass")

		log.Println("new pass = ", newPass)
		transfer.Send(Conn, newPass)

		log.Println("new pass send finish ", newPass)

		fmt.Fprint(w, "successfully change password to ", newPass, "<br>")
		BackHome(w)
	}
}
