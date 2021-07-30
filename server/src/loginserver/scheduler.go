package main

import (
	"errors"
	"sync"
)

/*
	Center Server Scheduler
 */

var (
	schedule *SchedulerMgr
	scheduleMutex sync.RWMutex
	ErrParseConfig = errors.New("Parse config error.")
)

func Scheduler_GetMe() (mgr *SchedulerMgr){
	scheduleMutex.RLock()
	if schedule == nil {
		schedule = &SchedulerMgr{
			rcenters: make(map[uint8]*CenterClient),
		}
	}
	mgr = schedule
	scheduleMutex.RUnlock()
	return
}

type SchedulerMgr struct {
	rcenters		map[uint8]*CenterClient

	configPath		string

}

///// Allocate room server	//////
func (this *SchedulerMgr) AllowFree(city uint32, cnet uint8) *CenterClient {
	client, ok := this.rcenters[cnet]
	if  ok && client.IsConnect() {
		return client
	}
	return client
}