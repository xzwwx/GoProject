package main

import (
	"common"
	"glog"
)

type RCenterClient struct {
	rate int
}

var rcenterm *RCenterClient

func RCenterClient_GetMe() *RCenterClient {
	if rcenterm == nil {
		rcenterm = &RCenterClient{}
	}
	return rcenterm
}

func (this *RCenterClient) GetRoom (userid uint64, isnew bool, others interface{} ) (*common.RetRoom, bool){
	center := Scheduler_GetMe().AllowFree(0,0)
	if center == nil {
		glog.Info("[Allocate] Get center server failed. ", 0, ", ", 0)
		return nil, false
	}
	reply := common.RetRoom{}
	err := center.RemoteCall("RPCTask.GetRoom", &common.ReqRoom{
		UserId: userid,
		//IsNew:  false,
	}, &reply)
	if err != nil {
		return nil, false
	}
	return &reply, true

}