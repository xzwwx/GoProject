package main

import (
	"fmt"
	"net/rpc"
	"time"
	message "usercmd"
)

func main(){

	client, err := rpc.DialHTTP("tcp", "localhost:8081")
	if err != nil {
		panic(err.Error())
	}

	timeStamp := time.Now().Unix()
	request := message.OrderRequest{OrderId: "201907310002", TimeStamp: timeStamp}

	var response *message.OrderInfo
	err = client.Call("OrderService.GetOrderInfo", request, &response)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(*response)
}
