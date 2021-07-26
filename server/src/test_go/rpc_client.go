package main
import (
	"fmt"
	go_protoc "usercmd"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply = &go_protoc.String{}
	var param = &go_protoc.String{
		Value:"hello",
	}

	err = client.Call("HelloService.Hello", &param, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}