package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type User struct{
	userName string
	userEmail string
	userPasswd string
	userTel string
	userKey string
	isNew bool
}

func password_md5(password, salt string)string{
	h := md5.New()
	h.Write([]byte(salt + password))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
		PoolSize: 100,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := client.Ping(ctx).Result()
	if err != nil{
		panic(err)
	}

	return client
}

func userLogin(client *redis.Client, userName , email, password, salt string)(isLogin, getName string){
	ctx := context.Background()
	if userName != ""{

		uidKey, err := client.Get(ctx, "userName:"+strings.ToLower(strings.TrimSpace(userName))).Result()
		if err == redis.Nil{
			return "userName does not exist", "00000000"
		}

		getpwd, err := client.HGet(ctx, uidKey, "password").Result()
		if err != nil{
			panic(err)
		}

		if getpwd == password_md5(password, salt){
			return "login sucessfully", userName
		}

	}else if email != ""{
		uidKey, err := client.Get(ctx, "email:"+email).Result()
		if err == redis.Nil{
			return "email does not exist", "00000000"
		}

		getpwd, err := client.HGet(ctx, uidKey, "password").Result()
		if err != nil{
			panic(err)
		}

		if getpwd == password_md5(password, salt){
			userName, err = client.HGet(ctx, uidKey, "userName").Result()
			if err != nil{
				panic(err)
			}
			return "login sucessfully", userName
		}

	}
	return "fail", "00000000"
}

// login POST
func login(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read request.Body failed, err:%v\n", err)
		return
	}
	data:= make(map[string]string)
	err = json.Unmarshal([]byte(body), &data)
	if err !=nil {
		fmt.Println("body To Json failed, err:", err)
		panic("Error")
	}
	fmt.Println(data)

	//login
	client := createClient()
	fmt.Println(client)

	var newResister User
	newResister.userName = data["userName"]
	newResister.userEmail = data["userEmail"]
	newResister.userPasswd = data["userPasswd"]
	newResister.userTel = data["userTel"]
	newResister.userKey = data["userKey"]

	status, userName := userLogin(client, newResister.userName, "", newResister.userPasswd, newResister.userName)
	if status == "login sucessfully" {
		fmt.Println(userName)
	}
	answer := `{"Login": "ok"}`
	w.Write([]byte(answer))

}

func main(){
	//client := createClient()
	//fmt.Println(client)
	//
	//register(client)
	//status, userName := userLogin(client,"oceanstar", "", "123456", "yan")
	//if status == "login sucessfully" {
	//	fmt.Println(userName)
	//}

	http.HandleFunc("/rigester", rigester)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil)
	log.Fatal(err)

}
