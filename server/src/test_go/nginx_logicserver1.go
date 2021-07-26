package main

import (
"fmt"
	//"net"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
	fmt.Fprintln(w, "r.Host = ", r.Host)
	fmt.Fprintln(w, "r.RequestURI = ", r.RequestURI)
}

func HelloWorldHandler(w http.ResponseWriter,r *http.Request)  {
	fmt.Println("r.Method = ", r.Method)
	fmt.Println("r.URL = ", r.URL)
	fmt.Println("r.Header = ", r.Header)
	fmt.Println("r.Body = ", r.Body)
	fmt.Fprintf(w,"HelloWorld!", r.Host)
}

func UserLoginHandler(response http.ResponseWriter,request *http.Request)  {
	fmt.Println("Handler Hello")
	fmt.Fprintf(response,"Login Success")
}

func main () {
	http.HandleFunc("/", HelloHandler)
	http.HandleFunc("/abc",HelloWorldHandler)
	http.HandleFunc("/user/login",UserLoginHandler)
	http.ListenAndServe(":9090", nil)
}