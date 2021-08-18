package main

import (
	"common"
	"fmt"
	"glog"
	"sync"
	"sync/atomic"
	"time"
	"usercmd"
)

// 初始数值
const (
	BombMaxTime = 4 // 炸弹的持续时间

	HurtScore = 1 // 造成伤害得分
	KillScore = 3 // 击杀得分

	RoleInitBombPower = 1  // 初始炸弹威力
	RoleInitSpeed     = 1  // 初始移动速度
	RoleInitHp        = 3  // 玩家初始血量
	RoleInitPosX      = 11 // 玩家初始位置x
	RoleInitPosY      = 15 // 玩家初始位置y
)

type ScenePlayer struct {
	id   uint64
	name string
	key  string // login key

	self  *PlayerTask
	scene *Scene

	otherPlayers map[uint64]*ScenePlayer
	rangeBombs   []*Bomb //current bombs

	senddie bool

	// neng fou fen li chu lai
	PlayerMove
	isMove    bool
	score     uint32           // 分数
	curPos    *common.Position // 当前位置
	nextPos   *common.Position // 下一个位置
	speed     float64
	power     uint32
	direction uint32
	lifeState uint32 // 0: dead   1: alive 	2: jelly
	hp        int32  // 生命值
	curbomb   uint32 // 当前已经放置的炸弹数量
	maxbomb   uint32 // 能放置的最大炸弹数量
	bombLeft  uint32

	//offline data
	udata *common.UserData

	movereq    *common.ReqMoveMsg
	laybombreq *common.ReqLayBombMsg
	objreq     *common.ReqTriggerObjectMsg // chi daoju
	killreq    *common.ReqKillMsg
}

func NewScenePlayer(player *PlayerTask, scene *Scene) *ScenePlayer {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &ScenePlayer{
		id:           player.id,
		name:         player.name,
		key:          player.key,
		scene:        scene,
		self:         player,
		otherPlayers: make(map[uint64]*ScenePlayer),
		senddie:      false,
		curPos:       &common.Position{X: RoleInitPosX, Y: float64(RoleInitPosY - (player.id%100)/1)},
		nextPos:      &common.Position{X: RoleInitPosX, Y: float64(RoleInitPosY - (player.id%100)/1)},
		hp:           RoleInitHp,
		score:        0,
		curbomb:      0,
		maxbomb:      5,
		power:        RoleInitBombPower,
		speed:        RoleInitSpeed,
		isMove:       false,
	}
	return s
}

//Sync from scene
// func (this *ScenePlayer) Update(perTime float64, scene *Scene) {
// 	//update move
// 	task := this.self
// 	if task != nil && atomic.CompareAndSwapInt32(&task.hasMove, 1, 0) {
// 		this.Move(scene, float64(atomic.LoadInt32(&task.speed)), float64(atomic.LoadInt32(&task.direction)))
// 	}

// 	//update bomb
// 	//this.rangebombs = this.rangebombs[:0]
// 	//this.rangebombs = append(this.rangebombs,this.bombs)
// 	for _, ball := range this.rangeBombs {
// 		if ball.isdelete == 1 {
// 			continue
// 		}
// 		scene.rangeBalls = append(scene.rangeBalls, ball)

// 		now := scene.now.Unix()
// 		if ball.layTime < now {
// 			continue
// 		}

// 	}

// }

// Send update msg to all
// func (this *ScenePlayer) sendSceneMsg(scene *Scene) {
// 	var (
// 		Moves       = scene.pool.moveFree[:0]
// 		ScenePlayer = scene.pool.msgPlayer[:0]
// 		MsgBomb     = scene.pool.msgBomb[:0]
// 		//Bombs =
// 	)
// 	// players move msg
// 	for _, playermove := range this.otherPlayers {
// 		if playermove.isMove == true {
// 			move := &usercmd.MsgPlayerMove{
// 				Id: playermove.id,
// 				X:  int32(playermove.pos.x),
// 				Y:  int32(playermove.pos.y),
// 				Nx: int32(playermove.nextpos.x),
// 				Ny: int32(playermove.nextpos.y),
// 			}
// 			Moves = append(Moves, move)
// 		}
// 		sp := &usercmd.ScenePlayer{
// 			Id:      playermove.id,
// 			BombNum: playermove.curbomb,
// 			Power:   playermove.power,
// 			Speed:   playermove.speed,
// 			State:   playermove.lifeState,
// 			X:       float32(playermove.pos.x),
// 			Y:       float32(playermove.pos.y),
// 			IsMove:  true,
// 		}
// 		ScenePlayer = append(ScenePlayer, sp)
// 	}

// 	for _, bomb := range this.rangeBombs {
// 		mb := &usercmd.MsgBomb{
// 			X:        bomb.pos.X,
// 			Y:        bomb.pos.Y,
// 			IsDelete: bomb.isdelete,
// 			Power:    bomb.power,
// 		}
// 		MsgBomb = append(MsgBomb, mb)
// 	}

// 	if len(Moves) != 0 {
// 		msg := &scene.pool.msgScene
// 		msg.Moves = Moves
// 		msg.Frame = scene.frame
// 		msg.Bombs = MsgBomb
// 		msg.Players = ScenePlayer
// 		msg.Id = this.id
// 		msg.X = float32(this.pos.x)
// 		msg.Y = float32(this.pos.y)

// 		this.sendSceneMsgToNet(msg, scene)
// 	}

// }

// 同步场景信息
func (this *ScenePlayer) sendSceneMsgToNet(msg *usercmd.MsgScene, scene *Scene) {
	if this.self != nil {
		// 1.通过msgEncode编码发送场景信息
		newPos := msgSceneToBytes(uint16(usercmd.MsgTypeCmd_SceneSync), msg, scene.msgBytes)
		//TCP
		this.self.AsyncSend(scene.msgBytes[:newPos], 0)

		// 2.通过 Grpc 自带方法编码 发送场景信息
		data, ok := common.EncodeToBytes(uint16(usercmd.MsgTypeCmd_SceneSync), msg)
		if !ok {
			glog.Info("[玩家]更新场景失败 cmd:", uint16(usercmd.MsgTypeCmd_SceneSync))
			return
		}
		this.self.AsyncSend(data, 0)

		// 3.通过Json编码的方式发送

		//data , _ = json.Marshal()

	}
}

// Update players in scene
func (this *ScenePlayer) UpdateViewPlayers(scene *Scene) {
}

func (this *ScenePlayer) AsyncSend(buffer []byte, flag byte) {

	this.self.AsyncSend(buffer, flag)
}

// Time Action
func (this *ScenePlayer) TimeAction(room *Room, timenow time.Time) bool {

	return true
}

// //Lay bomb
// func (this *ScenePlayer) LayBomb(room *Room, x, y int32) {

// 	var (
// 		isLayBomb = false
// 		scene     = &room.Scene
// 	)
// 	if this.bombLeft > 0 {
// 		// cell := scene.GetCellState(uint32(this.pos.x), uint32(this.pos.y))
// 		cell := scene.GetCellState(uint32(x), uint32(y))
// 		if cell == 0 { // ke fang zha dan
// 			timenow := time.Now().Unix()
// 			bomb := &Bomb{
// 				pos: VectorInt{
// 					// X: int32(this.pos.x),
// 					// Y: int32(this.pos.y),
// 					X: x,
// 					Y: y,
// 				},
// 				player:   this,
// 				layTime:  timenow,
// 				isdelete: 0,
// 			}
// 			this.bombLeft--

// 			// this.rangeBombs = append(this.rangeBombs, bomb)
// 			// scene.rangeBalls = append(scene.rangeBalls, bomb)

// 			// 另一种形式管理炸弹
// 			room.bombmgr.Add(bomb)

// 			isLayBomb = true
// 		}

// 	}
// 	if isLayBomb {
// 		this.SendCmd(usercmd.MsgTypeCmd_LayBomb, &usercmd.MsgLayBomb{})
// 	}

// 	if room != nil {

// 	}

// }

// Player Send Cmd
func (this *ScenePlayer) SendCmd(cmd usercmd.MsgTypeCmd, msg common.Message) bool {
	data, ok := common.EncodeToBytes(uint16(cmd), msg)
	if !ok {
		glog.Info("[Player] Send cmd:", cmd, ", len:", (len(data)))
		return false
	}
	this.AsyncSend(data, 0)
	return true
}

// 发送状态
func (this *ScenePlayer) SendState(room *Room) {
	room.BroadcastMsg(usercmd.MsgTypeCmd_PlayerState, PlayerStateCmd(this.id, int32(this.lifeState)))
}

//////////ScenePlayer manager//////////
type ScenePlayerMgr struct {
	mutex   sync.RWMutex
	players map[string]*ScenePlayer
}

var sptaskm *ScenePlayerMgr

func ScenePlayerMgr_GetMe() *ScenePlayerMgr {
	if sptaskm == nil {
		sptaskm = &ScenePlayerMgr{
			players: make(map[string]*ScenePlayer),
		}
	}
	return sptaskm
}

// Add
func (this *ScenePlayerMgr) Add(task *ScenePlayer) {
	this.mutex.RLock()
	this.players[task.key] = task
	this.mutex.RUnlock()
}

// Get player by key
func (this *ScenePlayerMgr) GetPlayer(key string) *ScenePlayer {
	this.mutex.RLock()
	player, _ := this.players[key]
	this.mutex.RUnlock()
	return player
}

// 删除场景玩家
func (this *ScenePlayerMgr) Removes(splayers map[uint64]*ScenePlayer) {
	this.mutex.Lock()
	for _, player := range splayers {
		delete(this.players, player.key)
	}
	fmt.Println("删除场景玩家")
	this.mutex.Unlock()
}

//////////////////////////////////////

//------lyf版本-----------------------------------------------------------------------------
// ----------------------玩家角色事件-------------------------- //

// 放置炸弹
func (this *ScenePlayer) PutBomb(msg *usercmd.MsgPutBomb) bool {
	// 达到最大炸弹数
	if atomic.LoadUint32(&this.curbomb) == this.maxbomb {
		return false
	}
	// 当前位置是否已经存在炸弹
	x, y := this.GetCurrentGrid()
	if this.scene.gameMap.MapArray[x][y] == GridType_Bomb {
		glog.Infof("[%v 放置炸弹] 位置{%v, %v}, 当前位置已存在炸弹",
			this.name, x, y)
		return false
	}

	bomb := NewBomb(this)
	this.scene.AddBomb(bomb)
	go bomb.CountDown()

	atomic.StoreUint32(&this.curbomb, this.curbomb+1)
	glog.Infof("[%v 放置炸弹] 炸弹位置{%v, %v}, 已放炸弹:%v, 剩余炸弹:%v",
		bomb.owner.name, bomb.pos.X, bomb.pos.Y, this.curbomb, this.maxbomb-this.curbomb)
	return true
}

// 移动
func (this *ScenePlayer) Move(msg *usercmd.MsgMove) {
	// TODO BUG:invalid memory address or nil pointer dereference
	if this == nil {
		return
	}
	this.isMove = true

	this.CaculateNext(msg.Way)     // 计算下一个位置
	this.BorderCheck(this.nextPos) // 边界检查

	this.curPos = this.nextPos
	glog.Infof("[%v 移动]当前位置为 x:%v, y:%v",
		this.name, this.nextPos.X, this.nextPos.Y)
}

// TODO 上下左右移动
func (this *ScenePlayer) CaculateNext(way int32) {
	x, y := this.GetCurrentGrid()
	// glog.Errorf("[GetCurrentGrid] x:%v, y:%v", x, y)
	// 上1下2左3右4
	switch common.MoveWay(way) {
	case common.MoveWay_Up:
		if this.CanPass(x, y+1) {
			this.nextPos.X = float64(x) // 如果在格子边缘，自动调整到格子中央
			this.nextPos.Y = this.curPos.Y + this.speed
		}
		break
	case common.MoveWay_Down:
		if this.CanPass(x, y-1) {
			this.nextPos.X = float64(x)
			this.nextPos.Y = this.curPos.Y - this.speed
		}
		break
	case common.MoveWay_Left:
		if this.CanPass(x-1, y) {
			this.nextPos.X = this.curPos.X - this.speed
			this.nextPos.Y = float64(y)
		}
		break
	case common.MoveWay_Right:
		if this.CanPass(x+1, y) {
			this.nextPos.X = this.curPos.X + this.speed
			this.nextPos.Y = float64(y)
		}
		break
	default:
	}
}

// 地图边界检查
func (this *ScenePlayer) BorderCheck(pos *common.Position) {
	this.scene.BorderCheck(pos)
}

// 该格子是否可以通过
func (this *ScenePlayer) CanPass(x, y uint32) bool {
	return this.scene.CanPass(x, y)
}

// 判断该玩家当前属于哪一个格子
func (this *ScenePlayer) GetCurrentGrid() (uint32, uint32) {
	return uint32(common.Round(this.curPos.X)),
		uint32(common.Round(this.curPos.Y))
}

// 收到伤害
func (this *ScenePlayer) BeHurt(attacker *ScenePlayer) {
	if atomic.StoreInt32(&this.hp, this.hp-1); this.hp <= 0 {
		// TODO 玩家角色死亡
		glog.Infoln("[玩家死亡] username:", this.name)
		info := &usercmd.RetRoleDeath{
			KillName: attacker.name,
			KillId:   attacker.id,
			LiveTime: 0,
			Score:    this.score,
		}
		// 向客户端发送玩家死亡信息
		this.self.SendCmd(usercmd.MsgTypeCmd_Death, info)
		// 玩家死亡房间/场景处理
		this.Death()
	}
	glog.Infof("[%v收到%v的炸弹的伤害] 当前血量hp:%v，%v当前得分:%v",
		this.name, attacker.name, this.hp, attacker.name, attacker.score)
}

// 玩家造成伤害或击杀，增加得分
func (this *ScenePlayer) AddScore(x uint32) {
	atomic.StoreUint32(&this.score, this.score+x)
	atomic.StoreUint32(&this.self.score, this.self.score+x)
}

// 发送场景同步信息
func (this *ScenePlayer) SendSceneMessage() {
	ret := &usercmd.RetUpdateSceneMsg{}
	// 场景内所有的玩家信息
	ret.Id = this.id
	for _, player := range this.scene.players {
		ret.Players = append(ret.Players, &usercmd.ScenePlayer{
			Id:      player.id,
			BombNum: player.curbomb,
			Power:   player.power,
			Speed:   float32(player.speed),
			State:   uint32(player.hp),
			X:       float32(player.curPos.X),
			Y:       float32(player.curPos.Y),
			IsMove:  player.isMove,
			Score:   this.score,
		})
	}
	// 场景内所有的炸弹信息
	for _, bomb := range this.scene.BombMap {
		ret.Bombs = append(ret.Bombs, &usercmd.MsgBomb{
			Id:         bomb.id,
			Own:        bomb.owner.id,
			X:          int32(bomb.pos.X),
			Y:          int32(bomb.pos.Y),
			CreateTime: 0,
		})
	}

	this.self.SendCmd(usercmd.MsgTypeCmd_SceneSync, ret)
}

func (this *ScenePlayer) Update() {

}

// 玩家死亡处理
func (this *ScenePlayer) Death() {
	this.scene.DelPlayer(this)             // 场景中删除
	this.self.room.RemovePlayer(this.self) // 把玩家从房间中删除

}

// 添加玩家数据到场景同步信息
func (this *ScenePlayer) AddAllPlayerInfoToMessage(msg common.Message) {
	// 类型断言
	if ret, ok := msg.(*usercmd.RetUpdateSceneMsg); ok {

		for _, player := range this.scene.players {
			ret.Players = append(ret.Players, &usercmd.ScenePlayer{
				Id:      player.id,
				BombNum: player.curbomb,
				Power:   player.power,
				Speed:   float32(player.speed),
				State:   uint32(player.hp),
				X:       float32(player.curPos.X),
				Y:       float32(player.curPos.Y),
				Score:   this.score,
				IsMove:  player.isMove,
			})
		}
	}
}

// ---------------------------------------------------------- //

///////////////原版本/////////////
// func (this *ScenePlayer) Move0(scene *Scene, speed, direction float64) {

// 	this.isMove = true
// 	this.CaculateNext(direction)          // 计算下一个位置
// 	this.scene.BorderCheck(&this.nextpos) // 保证计算得到的下一位置不超出地图范围

// 	// TODO 判断移动路径上是否有障碍物
// 	cellstate := scene.GetCellState(uint32(this.nextpos.x), uint32(this.nextpos.y))
// 	if cellstate == 0 {
// 		this.MoveVec(scene, speed, direction)
// 	} else if cellstate == 3 {
// 		// 吃道具
// 		this.MoveVec(scene, speed, direction)
// 		obj := scene.GetObjType(int32(this.nextpos.x), int32(this.nextpos.y))
// 		switch obj { // 1: 加速 speed   2：威力 power  3：数量  Bombnum
// 		case 1:
// 			this.speed++
// 			atomic.StoreInt32(&scene.gameMap.gamemap[int32(this.nextpos.x)][int32(this.nextpos.y)], 0)
// 		case 2:
// 			this.power++
// 			atomic.StoreInt32(&scene.gameMap.gamemap[int32(this.nextpos.x)][int32(this.nextpos.y)], 0)
// 		case 3:
// 			this.maxbomb++
// 			atomic.StoreInt32(&scene.gameMap.gamemap[int32(this.nextpos.x)][int32(this.nextpos.y)], 0)
// 		}
// 	} else {
// 		// 位置不改变
// 	}
// }

// func (this *ScenePlayer) MoveVec(scene *Scene, speed, direction float64) {
// 	this.pos = this.nextpos
// }

// // 计算下一个位置
// func (this *ScenePlayer) CaculateNext(direction float64) {
// 	this.nextpos.x = this.pos.x + float64(this.speed)*(math.Cos(direction*math.Pi/2))*0.04
// 	this.nextpos.y = this.pos.y + float64(this.speed)*(math.Sin(direction*math.Pi/2))*0.04

// }
