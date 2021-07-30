package main

import (
	"glog"
	"base/rpc"
	"sync/atomic"
	"time"
)

type CenterClient struct {
	client *rpc.Client
	isclose int32
	rcaddr string
	rate int
	zoneid uint8
	maxNum int
	curNum int
}

func NewRClient(addr string) *CenterClient {
	return &CenterClient{
		isclose: 1,
		rcaddr: addr,
	}
}

func (this *CenterClient) IsConnect() bool {
	return atomic.LoadInt32(&this.isclose) == 0
}

func (this *CenterClient) Connect() bool {
	if atomic.LoadInt32(&this.isclose) == 0 {
		return false
	}
	client, err := rpc.Dial("tcp", this.rcaddr, this.ReConnect)
	if err != nil {
		glog.Info("[RPC] Connect failed. ", err)
		return false
	}
	if client != nil {
		this.client = client
	}
	atomic.StoreInt32(&this.isclose, 0)
	glog.Info("[RPC] Connected to RCenterServer successfully. ", this.rcaddr)
	return true
}

func (this *CenterClient) ReConnect() bool {
	if atomic.CompareAndSwapInt32(&this.isclose, 0, 1) {
		for{
			glog.Info("[RPC] Reconnecting...")
			if this.Connect() {
				break
			}
			time.Sleep(time.Second * 2)
		}
		return true
	}
	return false
}

func (this *CenterClient) RemoteCallWithTT(serviceMethod string, args interface{}, reply interface{}, TT int) error {
	if this.client == nil {
		glog.Error("[RPC] Uninitialized. ", serviceMethod)
		return ErrDBClient
	}
	err := this.client.CallWithTimeout(serviceMethod, args, reply, time.Duration(TT) * time.Second)
	if err != nil {
		switch err{
		case rpc.ErrTimeout:
			this.client.Close()
			glog.Error("[RPC] Remote call timeout. ", this.rcaddr, ", ", serviceMethod)
			if this.ReConnect(){
				return this.client.CallWithTimeout(serviceMethod, args, reply, time.Duration(TT) * time.Second)
			}
		case rpc.ErrShutdown:
			if this.ReConnect(){
				return this.client.CallWithTimeout(serviceMethod, args, reply, time.Duration(TT) * time.Second)
			}
		}
		return err
	}
	return nil
}

func (this *CenterClient) RemoteCall (serviceMethod string, args interface{}, reply interface{}) error {
	timenow := time.Now()
	defer func(){
		usemsecs := time.Now().Sub(timenow).Milliseconds()
		glog.Info("[Kadun] ", this.rcaddr, ", ", serviceMethod, ", ", usemsecs, ", ", args)
	}()

	if this.client == nil {
		glog.Error("[RPC] Uninitialized. ", serviceMethod)
		return ErrDBClient
	}

	err := this.client.CallWithTimeout(serviceMethod, args, reply, 2 * time.Second)
	if err != nil {
		switch err{
		case rpc.ErrTimeout:
			this.client.Close()
			glog.Error("[RPC] Remote call timeout. ", this.rcaddr, ", ", serviceMethod)
			//if this.ReConnect(){
			//	return this.client.CallWithTimeout(serviceMethod, args, reply, time.Duration(TT) * time.Second)
			//}
		case rpc.ErrShutdown:
			if this.ReConnect(){
				return this.client.CallWithTimeout(serviceMethod, args, reply, 5 * time.Second)
			}
		}
		return err
	}
	return nil
}



