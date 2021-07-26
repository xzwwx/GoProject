package main

import (
	"fmt"
	"net/http"
)

type HelloHandlerStruct struct {
	content string
}

func (handler *HelloHandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, handler.content)
	fmt.Fprintln(w, "r.Host = ", r.Host)

}

func main()  {
	http.Handle("/", &HelloHandlerStruct{content: "Hello World"})
	http.ListenAndServe(":9091", nil)
}