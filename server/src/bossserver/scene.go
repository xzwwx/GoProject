package main

import (
	"common"
	"sync/atomic"
	"time"
	"usercmd"

	"glog"
)

const (
	RoomWidth  = 200
	RoomHeight = 200
	CellWidth  = 5
	CellHeight = 5

	RoomSize = 200
)

type Scene struct {
	room        *Room
	players     map[uint64]*ScenePlayer
	Obstacle    *map[uint32]*Obstacle
	ObstacleMap map[uint32]*common.Obstacle
	BombMap     map[uint32]*Bomb
	BoxMap      map[uint32]*common.Box
	gameMap     *GameMap // 游戏地图信息

	// Scene info
	sceneWidth     float64
	sceneHeight    float64
	now            time.Time
	startTime      time.Time
	frame          uint32
	pool           *BallPool // Player Pool
	msgBytes       []byte
	bombNum        uint32  // 炸弹编号
	rangeBalls     []*Bomb //map
	rangeObstacles []*Obstacle
}

func NewScene(room *Room) *Scene {
	scene := &Scene{
		room:        room,
		players:     make(map[uint64]*ScenePlayer),
		ObstacleMap: make(map[uint32]*common.Obstacle),
		BoxMap:      make(map[uint32]*common.Box),
		BombMap:     make(map[uint32]*Bomb),
		bombNum:     0,
		gameMap:     nil,
	}
	//scene.Init(room)
	return scene
}

func (this *Scene) Init(room *Room) {

	// 房间指针
	this.room = room

	this.players = make(map[uint64]*ScenePlayer)
	this.LoadGameMapData()

	this.rangeBalls = this.rangeBalls[:0]
	// this.rangePlayers = this.rangePlayers[:0]

	this.startTime = time.Now()

	// this.Obstacle = GenerateRandMap()

	// 加载地图
	// this.gameMap = &Map{} // 50x50

	// 初始化道具位置

}

// func (this *Scene) AddPlayer(p *PlayerTask) {
// 	this.players[p.id] = NewScenePlayer(p, this)
// }

// func (this *Scene) SendRoomMsg() {
// 	for _, p := range this.players {
// 		p.sendSceneMsg(this)
// 	}
// }

// func (this *Scene) UpdatePlayers(per float64) {

// 	// Depose player logic
// 	//if this.room.roomType
// 	//for i := 0; i < len(this.players); i++ {
// 	//	player, _ := this.players[]
// 	//}
// 	for _, player := range this.players {
// 		player.Update(per, this)
// 	}

// }

// // Check mapcell   0: null  1:wall  2:Obstacle 3: 道具
// func (this *Scene) GetCellState(x, y uint32) int32 {
// 	return this.gameMap.gamemap[x][y]
// }

// // Check mapcell   0: null  1:wall  2:Obstacle 3: 道具
// func (this *Scene) GetObjType(x, y int32) int32 {
// 	return this.objects.obj[x][y]
// }

// // 保证位置不超出地图范围
// func (this *Scene) BorderCheck(pos *Vector2) {
// 	if pos.x < 0 {
// 		pos.x = 0
// 	} else if pos.x >= this.sceneWidth {
// 		pos.x = this.sceneWidth - 0.01
// 	}
// 	if pos.y < 0 {
// 		pos.y = 0
// 	} else if pos.y >= this.sceneHeight {
// 		pos.y = this.sceneHeight - 0.01
// 	}
// }

// TODO 加载地图数据
func (this *Scene) LoadGameMapData() bool {
	if this.gameMap == nil {
		glog.Errorln("[Scene] load game map error")
		return false
	}
	this.gameMap.Height = uint32(len(this.gameMap.MapArray))
	this.gameMap.Width = uint32(len(this.gameMap.MapArray[0]))
	// 纵坐标优先遍历
	var x, y uint32
	for x = 0; x < this.gameMap.Height; x++ {
		for y = 0; y < this.gameMap.Width; y++ {

			idx := x*this.gameMap.Width + y // 二维转一维
			gridType := this.gameMap.MapArray[x][y]
			if gridType == GridType_Obstacle {
				this.ObstacleMap[idx] = &common.Obstacle{
					Id: idx,
					Pos: common.GridPos{
						X: x,
						Y: y,
					},
				}
			} else if gridType == GridType_Box {
				this.BoxMap[idx] = &common.Box{
					Id:    idx,
					Goods: 1, // TODO 宝箱里的物品
					Pos: common.GridPos{
						X: x,
						Y: y,
					},
				}
			}

		}
	}
	return true
}

// 自定义地图信息
func (this *Scene) RandGameMapData_AllSpace() {
	this.gameMap = &GameMap{}
	this.gameMap.Height, this.gameMap.Width = 40, 40
	var x, y, i uint32
	// 初始化
	this.gameMap.MapArray = make([][]GridType, this.gameMap.Height)
	for i = 0; i < this.gameMap.Height; i++ {
		this.gameMap.MapArray[i] = make([]GridType, this.gameMap.Width)
	}

	// 赋值
	glog.Errorln("[游戏地图初始化] 初始化开始")
	for x = 0; x < this.gameMap.Height; x++ {
		for y = 0; y < this.gameMap.Width; y++ {
			this.gameMap.MapArray[x][y] = GridType_Space
		}
	}
	glog.Errorln("[游戏地图初始化] 初始化完成")
}

func (this *Scene) Update() {
	// TODO
	// 	for _, player := range this.players {
	// 		//player.Update()
	// 	}
}

// 场景内添加一个玩家
func (this *Scene) AddPlayer(player *PlayerTask) {
	if player != nil {
		glog.Infoln("[场景添加玩家] username: ", player.name)
		sp := NewScenePlayer(player, this)
		this.players[player.id] = sp
		//player.scenePlayer = sp
	}
}

// 场景内删除一个玩家，但玩家只要不断开连接，依然会在房间中，以便结算
func (this *Scene) DelPlayer(player *ScenePlayer) {
	if player != nil {
		glog.Infoln("[场景删除玩家] username: ", player.name)
		delete(this.players, player.id)
		player = nil
	}
	// 当前场景内只剩一个玩家，游戏胜利，房间计算
	atomic.StoreUint32(&this.room.curPlayerNum, this.room.curPlayerNum-1)
	if this.room.curPlayerNum == 1 {
		this.GameSettle()
	}
	// 房间结束
	this.room.endchan <- true
}

// 添加一个炸弹
func (this *Scene) AddBomb(bomb *Bomb) {
	this.BombMap[bomb.id] = bomb
	this.gameMap.MapArray[bomb.pos.X][bomb.pos.Y] = GridType_Bomb
}

// 删除一个炸弹（炸弹爆炸）
func (this *Scene) DelBomb(bomb *Bomb) {
	delete(this.BombMap, bomb.id)
	this.gameMap.MapArray[bomb.pos.X][bomb.pos.Y] = GridType_Space
	bomb = nil
}

// 获取下一个炸弹的编号
func (this *Scene) GetNextBombId() uint32 {
	return atomic.AddUint32(&this.bombNum, 1)
}

// 根据坐标返回地图上对应格子的当前类型（空地，墙体）
func (this *Scene) GetGameMapGridType(x, y uint32) GridType {
	// glog.Errorf("[GetGameMapGridType] x:%v, y:%v", x, y)
	return this.gameMap.MapArray[x][y]
}

func (this *Scene) BorderCheck(pos *common.Position) {
	if pos.X < 0 {
		pos.X = 0
	} else if w := float64(this.gameMap.Width - 1); pos.X >= w {
		pos.X = w
	}
	if pos.Y < 0 {
		pos.Y = 0
	} else if h := float64(this.gameMap.Height - 1); pos.Y >= h {
		pos.Y = h
	}
}

func (this *Scene) CanPass(x, y uint32) bool {
	if x >= this.gameMap.Width || y >= this.gameMap.Height {
		return false
	}
	t := this.gameMap.GetGridByPos(x, y)
	return t != GridType_Obstacle && t != GridType_Box
}

// 定时发送场景信息，包括各个玩家的信息
func (this *Scene) SendRoomMessage() {
	for _, player := range this.players {
		player.SendSceneMessage()
	}
}

// 游戏结算
func (this *Scene) GameSettle() {
	// 场景内的最后一位玩家胜利
	for _, player := range this.players {
		glog.Infof("[游戏结束] winner:%v, 得分:%v",
			player.name, player.score)
	}
	rm := this.room
	rm.endTime = time.Now()
	// 计算游戏持续时间（s）
	rm.totalTime = uint64(rm.endTime.Sub(rm.startTime).Seconds())
	// 给房间内的所有玩家发送结算信息
	for _, player := range this.room.rplayers {
		glog.Infof("[玩家结算] username:%v, gametime:%vs, score:%v",
			player.name, rm.totalTime, player.score)
		ret := &usercmd.SettleMentInfo{}
		ret.Id = player.id
		ret.GameTime = rm.totalTime
		// 游戏结算时候，玩家信息里只包括id和分数
		for _, ptask := range this.room.rplayers {
			info := &usercmd.ScenePlayer{}
			info.Id = ptask.id
			info.Score = ptask.score
			ret.Players = append(ret.Players, info)
		}
		// 发送房间结束命令
		player.SendCmd(usercmd.MsgTypeCmd_EndRoom, ret)
	}
}
