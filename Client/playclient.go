package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main(){
	url := "http://127.0.0.1:8090"

	//json
	contentType := "application/json"

	data :=map[string]interface{}{
		"uid" : "10086",
		"username":"xzw",
	}

	//`{"uid":"10086", "username":"xzw"}`


	j, err := json.Marshal(data)
	if err != nil {
		fmt.Println(" Json Marshal failed, err :", err)
	}
	fmt.Println(string(j))

	res, err2 := http.Post(url, contentType, strings.NewReader(string(j)))
	fmt.Println("+++++++++++post ok")


	if err2 != nil{
		fmt.Printf("Post failed, err: %v \n", err2)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Get response failed, err: %v\n", err)
		return
	}
	fmt.Println(string(body))

}