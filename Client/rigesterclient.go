package main
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	//"io"
	//"net/http"
	//"strings"
)

type User struct{
	userName string
	userEmail string
	userPasswd string
	userTel string
	userKey string
	isNew bool
}

func main(){
	url := "http://127.0.0.1:80/rigester"
	//biaodan shuju
	//contentType := "application/x-www-form-urlencoded"
	//data := "name=itbsl&age=18"


	//json
	contentType := "application/json"

	//data :=`{"name":"xzw", "age":"26"}`

	//var user_data User
	rigester_data := map[string]interface{}{
		"userName" : "xzw40",
		"userPasswd":"123",
		"userEmail":"888840@qq.com",
		"userTel":"132456",
		"userKey":"6666",
	}
	j, err := json.Marshal(rigester_data)
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