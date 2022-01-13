package main

import (
	//"fmt"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main(){
	db, err := sql.Open("mysql", "yhchen2001:66Cyh90523@tcp(127.0.0.1:3306)/server")

    // if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
    }

    log.Println("connection success")

	insert, err := db.Query("INSERT INTO test VALUES ( 2, 'TEST' )")

    // if there is an error inserting, handle it
    if err != nil {
        panic(err.Error())
    }
    // be careful deferring Queries if you are using transactions
    defer insert.Close()

}