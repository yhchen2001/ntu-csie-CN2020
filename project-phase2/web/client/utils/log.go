package utils

import(
	"fmt"
	"log"
	"os"
)

func CheckErrMild(err error) {
	if err != nil{
		fmt.Fprintf(os.Stderr, "Mild error: %s\n", err)
	}
}

func CheckErrFatal(err error){
	if err != nil{
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
		os.Exit(1)
	}
}

func Log(v ...interface{}) {
	log.Println(v...)
}