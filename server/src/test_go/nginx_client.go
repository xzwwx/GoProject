package main


import (
	"io"
	"log"
	"net/http"
	"fmt"

)

func main(){
	res, err := http.Get("http://127.0.0.1:80")
	if err != nil {
		log.Fatalln("GET error:", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil{
		log.Fatalln("io.ReadAll err : ", err)

	}
	fmt.Println(string(body))
	//fmt.Println(string(body))
}


