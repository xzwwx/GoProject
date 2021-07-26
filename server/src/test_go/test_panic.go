package main

import (
	"fmt"
	"errors"
)

func main(){
	fmt.Println("Enter func")
	defer func(){
		fmt.Println("Enter defer")
		if p:=recover(); p!= nil{
			fmt.Println("panic:", p)
		}
		fmt.Println("Exit defer func")
	}()
	panic(errors.New("Something wrong"))
	fmt.Println("Exit func main")
}
