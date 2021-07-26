package main
import (
"encoding/json"
"fmt"
"io"
"net/http"
"strings"

)


func main(){
	url := "http://127.0.0.1:80/login"
	//biaodan shuju
	//contentType := "application/x-www-form-urlencoded"
	//data := "name=itbsl&age=18"

	//json
	contentType := "application/json"

	//data :=`{"name":"xzw", "age":"26"}`

	login_data := map[string]interface{}{
		"userName" : "xzw",
		"userPasswd":"123",
		"userEmail":"8888@qq.com",
		"userTel":"132456",
		"userKey":"6666",
	}
	j, err := json.Marshal(login_data)
	if err != nil {
		fmt.Println(" Json Marshal failed, err :", err)
	}
	fmt.Print(string(j))

	res, err := http.Post(url, contentType, strings.NewReader(string(j)))
	if err != nil{
		fmt.Printf("Post failed, err: %v \n", err)
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