package main

import (
	"common"
	"errors"
	"glog"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"usercmd"
)

type Room struct {
	//Scene
	mutex        sync.Mutex
	id           uint32                 //房间id
	roomType     uint32                 //房间类型
	rplayers     map[uint64]*PlayerTask //房间内的玩家
	curPlayerNum uint32                 //当前房间内玩家数
	bombCount    uint32

	isStart  bool
	timeLoop uint64
	stopCh   chan bool
	isStop   bool
	iscustom bool
	frame    uint32
	isclosed int32

	rmode RoomMode

	// 处理炸弹
	//bombmgr *BombMgr

	// Operation
	scene             *Scene
	opChan            chan *opMsg // player operation msg
	chan_PlayerOp     chan *PlayerOp
	chan_Control      chan int
	chan_AddPlayer    chan *PlayerTask
	chan_RemovePlayer chan *PlayerTask

	newLogin      map[uint64]bool
	loginHis      map[uint64]bool
	newLoginMutex sync.Mutex

	preFrameTime time.Time

	msgBytes []byte

	now          time.Time
	startTime    time.Time
	endTime      time.Time // 结束时间
	maxPlayerNum uint32    //max player number. default :8
	totalTime    uint64    //in second
	endchan      chan bool
}

type opMsg struct {
	op   uint32
	args interface{}
}

type RoomMode interface {
	LoadData() bool
}

//type PlayerOpt struct {
//	pTask 		*PlayerTask
//	pPlayer 	*ScenePlayer
//	Opts		*UserOpt
//}

func NewRoom(rtype, rid uint32, player *PlayerTask) *Room {
	glog.Errorln("New房间：", rid, "开始")

	room := &Room{
		id:             rid,
		roomType:       rtype,
		rplayers:       make(map[uint64]*PlayerTask),
		curPlayerNum:   0,
		maxPlayerNum:   2,
		isStart:        false,
		isStop:         false,
		isclosed:       -1,
		endchan:        make(chan bool),
		chan_PlayerOp:  make(chan *PlayerOp, 500),
		chan_AddPlayer: make(chan *PlayerTask, 10),
	}
	room.scene = NewScene(room) // 初始化场景信息

	glog.Errorln("New房间：", rid, "完毕")

	return room
}

func (this *Room) Start() bool {
	if !atomic.CompareAndSwapInt32(&this.isclosed, -1, 0) {
		return false
	}
	//this.rmode = NewFreeRoom(this)
	//this.isStart = true

	//this.bombmgr = NewBombMgr(this)
	this.scene.Init(this)

	go this.Loop()
	glog.Info("[房间] 创建房间 ", ", ", this.id, ", ")
	return true
}

// 房间停止
func (this *Room) Stop() bool {
	if !atomic.CompareAndSwapInt32(&this.isclosed, 0, 1) {
		return false
	}
	this.destory()
	glog.Info("[房间] 销毁房间 ", this.id, ", ", len(this.rplayers))

	return true
}

func (this *Room) IsClosed() bool {
	return atomic.LoadInt32(&this.isclosed) != 0
}

// 删除房间
func (this *Room) destory() {
	this.Stop()
	go func(room *Room) {
		ScenePlayerMgr_GetMe().Removes(room.scene.players)

		// redis 清理玩家
		// TODO

	}(this)
	glog.Info("[房间] 结算完成", this.id, ", ", this.GetPlayerNum())
}

func (this *Room) AddLoginUser(UID uint64) (result bool) {
	this.newLoginMutex.Lock()
	defer this.newLoginMutex.Unlock()

	result = this.loginHis[UID]
	if !result {
		this.loginHis[UID] = true
	}
	this.newLogin[UID] = true
	return
}

// 玩家 +1
func (this *Room) IncPlayerNum() {
	atomic.AddUint32(&this.curPlayerNum, 1)
}

// 返回玩家数
func (this *Room) GetPlayerNum() int32 {
	return int32(atomic.LoadUint32(&this.curPlayerNum))

}

// Main game loop
func (this *Room) Loop() {

	this.scene.gameMap = &GameMap{}
	// TODO 加载地图信息
	if !this.scene.gameMap.CustomizeInit() {
		glog.Errorln("[地图加载失败]")
		return
	}

	timeTicker := time.NewTicker(time.Millisecond * 20)
	//stop := false
	defer func() {
		this.Stop()
		timeTicker.Stop()
		RoomMgr_GetMe().RemoveRoom(this)
	}()
	stop := false

	for !stop {
		this.now = time.Now()
		select {
		case <-timeTicker.C:
			// 0.02s
			if this.isStart == true {
				if this.timeLoop%2 == 0 {
					//this.Update(0.04)
				}
				//0.1s
				if this.timeLoop%5 == 0 {
					this.frame++
					this.scene.SendRoomMessage()
					//this.SendRoomMsg()
					//this.bombmgr.ExecAction()
				}

				//1s
				// if this.timeLoop%100 == 0 {

				// 	//this.scene.sendTime(this.totalTime - this.timeLoop/100)
				// }
				// if this.timeLoop != 0 && this.timeLoop%(this.totalTime*100) == 0 {
				// 	// stop = true
				// }
				this.timeLoop++

			}

			if this.isStop {
				stop = true
			}
		case op := <-this.chan_PlayerOp:
			//this.scene.UpdateOP(op)
			if this.isStart == true {
				switch op.opType {
				case PlayerLayBombOp:
					this.LayBomb(op.playerId, op.x, op.y)
				case PlayerMoveOp:
					req, ok := op.msg.(*usercmd.MsgMove)
					if !ok {
						glog.Errorln("[Move] move arg error")
						return
					}
					this.scene.players[op.playerId].Move(req)
				}
			}
		case player := <-this.chan_AddPlayer:
			glog.Errorln("----------------添加玩家chan")
			this.AddPlayer(player)
		}

	}
	//stop = true
	this.Close()
}

func (this *Room) Close() {
	if !this.isStop {
		//this.scene.SendOverMsg()
		this.isStop = true
		RoomMgr_GetMe().endChan <- this.id
	}
}

func (this *Room) Update(per float64) {
	//this.scene.UpdatePos()
	starttime := time.Now()
	ftime := starttime.Sub(this.preFrameTime).Milliseconds()
	this.preFrameTime = starttime

	//this.UpdatePlayers(per)
	rtime := time.Now().Sub(starttime).Milliseconds()
	if math.Abs(float64(ftime-40)) > 20 || rtime > 20 {
		glog.Info("[Statistic] State sync.", this.roomType, ", ", this.id, ", ", this.frame, ", ", ftime, ", ", rtime)
	}
}

//广播消息
func (this *Room) BroadcastMsg(msgNo usercmd.MsgTypeCmd, msg common.Message) {
	data, ok := common.EncodeToBytes(uint16(msgNo), msg)
	if !ok {
		glog.Error("[广播] 发送消息失败 ", msgNo)
		return
	}
	for _, player := range this.scene.players {
		player.AsyncSend(data, 0)
	}
}

func (this *Room) TimeAction() {

}

//Send cmd to room pthread
func (this *Room) Control(ctrl int) bool {
	if this.IsClosed() {
		return false
	}
	this.chan_Control <- ctrl
	return true

}

//Lay bomb
func (this *Room) LayBomb(playerId uint64, x, y int32) {
	// player, ok := this.players[playerId]
	// if !ok {
	// 	return
	// }
	//player.LayBomb(this, x, y)

}

// 玩家进入房间
func (this *Room) AddPlayer(player *PlayerTask) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.curPlayerNum > this.maxPlayerNum {
		return errors.New("room is full")
	}
	// 更新房间信息
	// this.curPlayerNum++
	player.room = this
	this.rplayers[player.id] = player
	glog.Errorln("[房间] 玩家[%v]进入[%v]房间  ", player.name, this.id)
	this.scene.AddPlayer(player) // 将玩家添加到场景

	// 房间内玩家数量达到最大，自动开始游戏
	if this.curPlayerNum == this.maxPlayerNum {
		glog.Infoln("[房间] 玩家数量：", len(this.rplayers))
		glog.Infoln("[游戏开始] 玩家列表：")
		for _, v := range this.rplayers {
			glog.Infof("username:%v, uid:%v", v.name, v.id)
		}
		RoomMgr_GetMe().UpdateNextRoomId() // 房间id++

		// 将当前房间内的所有玩家信息发送到客户端
		for _, pt := range this.rplayers {
			info := &usercmd.RetUpdateSceneMsg{}
			info.Id = pt.id
			for _, sp := range this.scene.players {
				info.Players = append(info.Players, &usercmd.ScenePlayer{
					Id:      sp.id,
					BombNum: sp.curbomb,
					Power:   sp.power,
					Speed:   float32(sp.speed),
					State:   uint32(sp.hp),
					X:       float32(sp.pos.x),
					Y:       float32(sp.pos.y),

					IsMove: sp.isMove,
				})
			}
			pt.SendCmd(usercmd.MsgTypeCmd_SceneSync, info)
		}

		this.isStart = true
		//go this.Start()
	}

	return nil
}

// 将玩家移除出房间
func (this *Room) RemovePlayer(player *PlayerTask) error {
	if this == nil {
		return nil
	}
	this.mutex.Lock()
	glog.Warningln("[debug]Room.RemovePlayer() func")
	defer this.mutex.Unlock()
	delete(this.rplayers, player.id)
	glog.Warningln("[debug]Room.RemovePlayer() func")
	return nil
}
