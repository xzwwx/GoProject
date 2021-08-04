package main

import (
	"base/gonet"
	"glog"
	"common"
	"net"
	"sync"
	"time"
)

const (
	Task_Max_Timeout = 1
	OpsPerSecond = 12
)

type PlayerTask struct {
	tcptask *gonet.TcpTask
	//udptask *snet.Session
	isUdp bool

	key string
	id uint64
	name string
	room *Room
	udata *common.UserData
	uobjs []uint32

	power uint32
	speed uint32
	lifenum uint32
	state uint32

	activeTime	time.Time
	onlineTime	int64

}

type PlayerOpType int

const (
	PlayerNoneOp 	= PlayerOpType(iota)
	PlayerMoveOp
	PlayerLayBombOp
	PlayerCombineOp
)

type PlayerOp struct {
	player  	*PlayerTask
	cmdParam	uint32
	opType 		PlayerOpType
	loginUsers	map[uint64]bool
	toPlayerId	uint64
	Opts 		*UserOpt
}

func NewPlayerTask(conn net.Conn) *PlayerTask{
	s := &PlayerTask{
		tcptask:    gonet.NewTcpTask(conn),
		activeTime: time.Now(),
		onlineTime: time.Now().Unix(),
		isUdp:      true,
	}
	s.tcptask.Derived = s
	return s
}

func (this *PlayerTask) ParseMsg(data []byte, flag byte) bool {



	return true
}

func (this *PlayerTask) OnClose() {

}

func (this *PlayerTask) Start(){
	if !this.isUdp {
		this.tcptask.Start()
	}
}

func (this *PlayerTask) Stop() bool {
	if this.isUdp {
		return true
	}else{
		this.tcptask.Close()
	}
	return true
}






//////////////////PlayerTask Manager//////////
type PlayerTaskMgr struct {
	mutex sync.RWMutex
	tasks map[uint64]*PlayerTask
}

var ptaskm *PlayerTaskMgr

func PlayerTaskMgr_GetMe() *PlayerTaskMgr {
	if ptaskm == nil {
		ptaskm = &PlayerTaskMgr{
			tasks: make(map[uint64]*PlayerTask),
		}
		go ptaskm.timeAction()
	}
	return ptaskm
}

// Depose player timeout
func (this *PlayerTaskMgr) timeAction(){


}

// Add
func (this *PlayerTaskMgr) Add(task *PlayerTask) bool {
	if task == nil {
		return false
	}
	this.mutex.Lock()
	this.tasks[task.id] = task
	this.mutex.Unlock()
	return true
}

func (this *PlayerTaskMgr) Remove(task *PlayerTask) bool {
	if task == nil {
		return false
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	t, ok := this.tasks[task.id]
	if !ok {
		return false
	}
	if t != task {
		glog.Error("[Logout] Failed. ", t.id, ", ", &t, ", ", &task)
		return false
	}

	delete(this.tasks, task.id)

	return true
}

func (this *PlayerTaskMgr) GetTask(uid uint64) *PlayerTask {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	user, ok := this.tasks[uid]
	if !ok {
		return nil
	}
	return user
}

func (this *PlayerTaskMgr) GetNum() int32 {
	this.mutex.RLock()
	tasknum := int32(len(this.tasks))
	this.mutex.RUnlock()
	return tasknum
}