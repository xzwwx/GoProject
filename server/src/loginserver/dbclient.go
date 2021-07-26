package main

import (
	"errors"
	"net/rpc"
)

var ErrDBClient = errors.New("DB error.")

type DbClient struct {
		client *rpc.Client
		dbaddr string
		isclose bool
		errnum int
}

func (this *DbClient) RemoteCall(serviceMethod string, args interface{}, reply interface{}) error{

	err := errors.New("err")
	return err
}