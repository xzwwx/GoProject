package main

import (
	"base/gonet"
	"common"
	"encoding/json"
	"glog"
	"net"
	"sync"
	"time"
	"usercmd"
)

const (
	Task_Max_Timeout = 1
	OpsPerSecond     = 3  // max lay bomb per second: zui duo mei miao fang 3 ge bombs
	OpsNumPerSecond  = 10 // mei miao zui duo cao zuo 10 ci
)

type PlayerTask struct {
	tcptask *gonet.TcpTask
	//udptask	 *snet.Session
	isUdp bool

	key   string
	id    uint64
	name  string
	room  *Room
	udata *common.UserData
	uobjs []uint32

	direction int32 // direction
	power     int32
	speed     int32 //Speed
	lifenum   uint32
	state     uint32
	hasMove   int32
	score     uint32 // 得分

	lastLayBomb     int32 // last lay bomb times
	lastLayBombTime int64 // last .. time

	activeTime time.Time
	onlineTime int64
}

type PlayerOpType int

const (
	PlayerNoneOp    = PlayerOpType(iota)
	PlayerMoveOp    // 1: 移动
	PlayerLayBombOp // 2: 放炸弹
	PlayerCombineOp
	PlayerEatObject
)

type PlayerOp_FrameSync struct {
	player     *PlayerTask
	cmdParam   uint32
	opType     PlayerOpType
	loginUsers map[uint64]bool
	toPlayerId uint64
	Opts       *UserOpt
}

type PlayerOp struct {
	playerId   uint64
	cmdParam   uint32
	opType     PlayerOpType
	loginUsers map[uint64]bool
	toPlayerId uint64
	opTime     uint64
	x          int32
	y          int32
	msg        common.Message // 其他信息
}

func NewPlayerTask(conn net.Conn) *PlayerTask {
	s := &PlayerTask{
		tcptask:    gonet.NewTcpTask(conn),
		activeTime: time.Now(),
		onlineTime: time.Now().Unix(),
		isUdp:      false,
	}
	s.tcptask.Derived = s
	return s
}

func (this *PlayerTask) ParseMsg(data []byte, flag byte) bool {
	this.activeTime = time.Now()

	if len(data) < 2 {
		return true
	}

	info := usercmd.CmdHeader{}
	err := json.Unmarshal(data, &info)
	if err != nil {
		glog.Errorln("[json解析失败] ", err)
	}
	cmd := info.Cmd

	//cmd := usercmd.MsgTypeCmd(common.GetCmd(data))
	if !this.IsVerified() {
		glog.Infoln("[debug] cmd : ", cmd)

		// Verify login
		if cmd != usercmd.MsgTypeCmd_Login1 {
			glog.Error("[Login] Not login cmd. ", this.RemoteAddr(), ", ", cmd)
			return false
		}

		// 解析Json消息
		revCmd := &usercmd.UserLoginInfo{}
		err := json.Unmarshal([]byte(info.Data), revCmd)
		if err != nil {
			glog.Errorln(err)
			// this.retErrorMsg(common.ErrorCodeRoom)
			return false
		}
		glog.Infoln("[RoomServer Login] recv a login request ", this.RemoteAddr())
		glog.Infoln(revCmd.Token)

		// ------原proto格式命令-----
		// revCmd, ok := common.DecodeGprotoCmd(data, flag, &usercmd.MsgLogin{}).(*usercmd.MsgLogin)
		// if !ok {
		// 	// return error msg
		// 	// SendCmd()
		// }
		//glog.Info("[Login] Received login request", this.RemoteAddr(), ", ", revCmd.Key, ", ")

		// 解析token，获取用户信息
		rinfo := &common.RoomTokenInfo{}
		glog.Errorln("Token:", revCmd.Token)

		// rinfo : userid, username, roomid
		rinfo, err = common.ParseRoomToken(revCmd.Token)
		glog.Errorln(rinfo.RoomId, ", ", rinfo.UserId, ", ", rinfo.UserName)
		if err != nil {
			glog.Errorln("[MsgTypeCmd_Login] parse room token error:", err)
			return false
		}
		// 同一个用户重复连接
		//遍历server里的所有房间
		for _, rm := range RoomMgr_GetMe().runrooms {
			for _, pl := range rm.rplayers {
				if pl.name == rinfo.UserName {
					this.OnClose()
					return false
				}
			}
		}

		//Check Key
		var newLogin bool = true
		if s := ScenePlayerMgr_GetMe().GetPlayer(rinfo.UserName); s != nil {
			this.udata = s.udata
			newLogin = false
		}
		if this.udata == nil {
			// udata := &common.UserData{}
			token := common.RedisMgr.Get(rinfo.UserName + "_roomtoken")
			if token == "" {
				glog.Error("[登录] 验证失败", this.RemoteAddr(), ", ", rinfo.UserName, "====")
				return false
			}
			udata, err := common.ParseRoomToken(token)
			glog.Errorln("Token RoomId:", udata.RoomId, ", UserId:", udata.UserId, ", Username:", udata.UserName)
			if err != nil {
				glog.Error("[登录] 验证失败", this.RemoteAddr(), ", ", rinfo.UserName)
				return false
			}
			ud := &common.UserData{
				Id:      udata.UserId,
				Account: udata.UserName,
				RoomId:  udata.RoomId,
			}
			this.udata = ud
			glog.Errorln("=======roomid:", this.udata.RoomId, ", id:", this.udata.Id)
			//this.udata.Id = udata.UserId
			//this.udata.Account = udata.UserName
			//this.udata.RoomId = udata.RoomId

			// if !RedisMgr_GetMe().LoadFromRedis(revCmd.Key, udata){
			// 	glog.Error("[登录] 验证失败", this.RemoteAddr(), ", ", revCmd.Key)
			// 	return false
			// }
		}

		this.key = rinfo.UserName
		this.id = this.udata.Id
		this.name = this.udata.Account

		otask := PlayerTaskMgr_GetMe().GetTask(this.id)
		if otask != nil {
			glog.Infoln("[Login] ReLogin.", otask.id, ", ", otask.key)
		}
		this.Verify()

		PlayerTaskMgr_GetMe().Add(this)

		if newLogin {
			room := RoomMgr_GetMe().getRoomById(this.udata.RoomId)
			if room == nil {
				room = RoomMgr_GetMe().NewRoom(0, rinfo.RoomId, this)
			}
		}

		//room := RoomMgr_GetMe().getRoomById(this.udata.RoomId)
		// if room != nil {
		// 	// 重连
		// }

		glog.Info("[Login] Verified account success. ", this.RemoteAddr(), ", ", this.udata.Id, ", ", this.udata.Account, ", ", this.key)

		//var joinroomtype uint32

		//if this.udata.Model

		if !RoomMgr_GetMe().AddPlayer(this) {
			return false
		}

		//this.online()

		glog.Info("[Login] Success,", this.RemoteAddr(), ", ", this.udata.RoomAddr, ", ", this.udata.RoomId, ", ",
			this.udata.Id, ", ", this.udata.Account, ", ", this.key)
		return true
	}

	//heartbeat
	//if cmd == usercmd.MsgTypeCmd_HeartBeat1 {
	//
	//}

	if this.room == nil || this.room.IsClosed() {
		glog.Info("[Message] Room end.")
		return false
	}

	switch cmd {
	case usercmd.MsgTypeCmd_Move:
		// Player move
		revCmd := &usercmd.MsgMove{}
		json.Unmarshal([]byte(info.Data), revCmd)
		if this.room == nil || this.room.IsClosed() {
			glog.Infoln("[收到请求移动的指令] 房间不存在")
			return false
		}
		if !this.room.isStart { // 游戏未开始
			glog.Infoln("[收到请求移动的指令] 游戏未开始")
			return false
		}
		glog.Infof("[%v收到请求移动的指令] revCmd.Way=%v", this.name, revCmd.Way)
		this.room.chan_PlayerOp <- &PlayerOp{playerId: this.id, opType: PlayerMoveOp, msg: revCmd}

		// if common.DecodeGoCmd(data, flag, revCmd) != nil {
		// 	return false
		// }

		// if revCmd.Speed == 0 {
		// 	revCmd.Speed = this.speed
		// }

		// atomic.StoreInt32(&this.direction, revCmd.Direction)
		// atomic.StoreInt32(&this.speed, revCmd.Speed)
		// atomic.StoreInt32(&this.hasMove, 1)

		// this.room.chan_PlayerOp <- &PlayerOp{playerId: this.id, cmdParam: 0, opType: PlayerMoveOp}
		//fmt.Println("Move++", this.id)

	case usercmd.MsgTypeCmd_LayBomb:
		// Lay bombs

		revCmd := &usercmd.MsgPutBomb{}
		json.Unmarshal([]byte(info.Data), revCmd)
		// if common.DecodeGoCmd(data, flag, revCmd) != nil {
		// 	return false
		// }
		if this.room == nil || this.room.IsClosed() {
			glog.Infoln("[收到请求放炸弹的指令] 房间不存在")
			return false
		}
		if !this.room.isStart { // 游戏未开始
			glog.Infoln("[收到请求放炸弹的指令] 游戏未开始")
			return false
		}
		glog.Infof("[%v收到请求放炸弹的指令]", this.name)
		this.room.chan_PlayerOp <- &PlayerOp{playerId: this.id, opType: PlayerLayBombOp, msg: revCmd}

		// timenow := time.Now().Unix()
		// if this.lastLayBombTime <= timenow {
		// 	if this.lastLayBomb >= OpsNumPerSecond {
		// 		glog.Error("[Lay Bomb] Too fast. ", this.udata.RoomId, ", ", this.udata.Id, ", ", this.udata.Account, ", ", this.lastLayBomb)
		// 	}
		// 	this.lastLayBombTime = timenow + 1
		// 	this.lastLayBomb = 0
		// }
		// this.lastLayBomb++
		// if this.lastLayBomb > OpsPerSecond {
		// 	return true
		// }
		// revCmd := &usercmd.MsgLayBomb{}
		// if common.DecodeGoCmd(data, flag, revCmd) != nil {
		// 	return false
		// }

		// this.room.chan_PlayerOp <- &PlayerOp{playerId: this.id, cmdParam: 0, opType: PlayerLayBombOp, opTime: uint64(revCmd.LayTime), x: int32(revCmd.X), y: int32(revCmd.Y)}

		// fmt.Println("Lay Bomb++", this.id)

	case usercmd.MsgTypeCmd_Death:
		//

	case usercmd.MsgTypeCmd_BeBomb:
		// Be bombed: bei zha dao

	case usercmd.MsgTypeCmd_EatObject:
		// Eat object
	case usercmd.MsgTypeCmd_Combine:
		//

	default:
		glog.Error("[Player] Unknown Cmd. ", this.id, ", ", this.name, ", ", cmd)
	}
	return true
}

func (this *PlayerTask) Verify() {
	this.tcptask.Verify()
}
func (this *PlayerTask) IsVerified() bool {
	// if this.isUdp
	return this.tcptask.IsVerified()
}

func (this *PlayerTask) OnClose() {
	if !this.IsVerified() {
		return
	}
	// offline delete from room

}

func (this *PlayerTask) RemoteAddr() string {
	return this.tcptask.RemoteAddr()
}

func (this *PlayerTask) Start() {
	if !this.isUdp {
		this.tcptask.Start()
	}
}

func (this *PlayerTask) Stop() bool {
	if this.isUdp {
		return true
	} else {
		this.tcptask.Close()
	}
	return true
}

//player online ; refresh room and server
func (this *PlayerTask) online() {
	room := this.room
	if room != nil && !room.IsClosed() {
		// update
		// To do
		// RedisMgr_GetMe()

		// RCenterClient_GetMe().UpdateRoom()

		go func() {
			//deng lu li shi
			room.AddLoginUser(this.id)

		}()

	}
	RCenterClient_GetMe().UpdateServer(RoomMgr_GetMe().getNum(), PlayerTaskMgr_GetMe().GetNum())
}

func (this *PlayerTask) offline() {

}

func (this *PlayerTask) SendCmd(cmd usercmd.MsgTypeCmd, msg common.Message) {
	// data, ok := common.EncodeToBytes(uint16(cmd), msg)
	data, ok := common.EncodeToBytesJson(uint16(cmd), msg)

	if !ok {
		glog.Info("[玩家] 发送消息失败 cmd:", cmd)
		return
	}
	this.AsyncSend(data, 0)
}

func (this *PlayerTask) AsyncSend(buffer []byte, flag byte) bool {

	return this.tcptask.AsyncSend(buffer, flag)
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
func (this *PlayerTaskMgr) timeAction() {

}

// Add
func (this *PlayerTaskMgr) Add(task *PlayerTask) bool {
	if task == nil {
		return false
	}
	this.mutex.Lock()
	this.tasks[task.id] = task
	glog.Errorln("添加playerTask：", task.id)
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
