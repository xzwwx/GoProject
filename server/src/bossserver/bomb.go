package main

import (
	"common"
	"time"

	"glog"
)

type Bomb struct {
	id    uint32          // 炸弹id，主要用于做Map的key
	pos   *common.GridPos // 位置
	owner *ScenePlayer    // 所有者
	scene *Scene          // 场景指针
}

func NewBomb(player *ScenePlayer) *Bomb {
	row, col := player.GetCurrentGrid()
	bomb := &Bomb{
		id:    row*player.scene.gameMap.Width + col,
		pos:   &common.GridPos{X: row, Y: col},
		owner: player,
		scene: player.scene,
	}
	// go func() {
	// 	ticker := time.NewTicker(BOMB_MAXTIME * time.Second)
	// 	<-ticker.C

	// 	bomb.Explode()

	// 	return
	// }()
	return bomb
}

// 倒计时
func (this *Bomb) CountDown() {
	// ticker := time.NewTicker(BOMB_MAXTIME * time.Second)
	// <-ticker.C
	time.Sleep(5 * time.Second)
	this.Explode()
}

// 爆炸
func (this *Bomb) Explode() {
	glog.Infof("[%v的炸弹爆炸] x:%v, y:%v",
		this.owner.name, this.pos.X, this.pos.Y)
	// 计算伤害范围
	// 1. 上下左右
	up := this.pos.Y + this.owner.power
	down := this.pos.Y - this.owner.power
	left := this.pos.X - this.owner.power
	right := this.pos.X + this.owner.power
	// 遍历所有炸弹，判断是否在当前炸弹的范围内(一颗炸弹引爆另一颗炸弹)
	for _, b := range this.scene.BombMap {
		if b.pos.Y == this.pos.Y && left <= b.pos.X && b.pos.X <= right {
			// b.Explode()
		}
		if b.pos.X == this.pos.X && down <= b.pos.Y && b.pos.Y <= up {
			// b.Explode()
		}
	}
	// 遍历所有角色，判断是否在当前炸弹的范围内
	for _, p := range this.scene.players {

		x, y := p.GetCurrentGrid()
		glog.Infof("[%v的炸弹爆炸时%v位置] x:%v, y:%v",
			this.owner.name, p.name, x, y)
		if y == this.pos.Y && left <= x && x <= right {
			// 水平方向
			glog.Infof("[炸弹造成<%v,%v>水平方向伤害]", left, right)
			this.owner.AddScore(HurtScore)
			p.BeHurt(this.owner)
		} else if x == this.pos.X && down <= y && y <= up {
			// 垂直方向
			glog.Infof("[炸弹造成<%v,%v>垂直方向伤害]", down, up)
			this.owner.AddScore(HurtScore)
			p.BeHurt(this.owner)
		}
	}
	// 在场景中删除炸弹
	this.scene.DelBomb(this)
	//
	this.owner.curbomb--
}

// package main

// import (
// 	"time"
// 	"usercmd"
// )

// type BaseBomb interface {
// 	Init(room *Room) bool
// 	Check(room *Room) bool
// 	Exec(room *Room) bool
// }

// type Bomb struct {
// 	id               int64
// 	player           *ScenePlayer
// 	layTime          int64
// 	pos              VectorInt
// 	power            int32
// 	lastTime         int64 // 爆炸时间
// 	isdelete         int32 // zha le ma
// 	boxDistroiedList []*RetMapCellState
// }

// func (this *Bomb) Init(room *Room) bool {
// 	timenow := room.now.Unix()
// 	this.lastTime = timenow + time.Second.Milliseconds()*3 // 3秒后爆炸
// 	this.pos.X = int32(this.player.pos.x)
// 	this.pos.Y = int32(this.player.pos.y)
// 	this.power = this.player.self.power
// 	room.rangeBalls = append(room.rangeBalls, this)
// 	this.isdelete = 0

// 	return true
// }

// func (this *Bomb) Check(room *Room) bool {
// 	timenow := time.Now().Unix()
// 	if timenow > this.lastTime { // 到爆炸时间了
// 		return true
// 	}
// 	return false

// }

// func (this *Bomb) Exec(room *Room) bool {
// 	// 先放的先炸
// 	bomb := room.rangeBalls[0]

// 	if bomb.layTime == this.layTime {

// 		isbombwalls := this.IsBombWalls(room)
// 		if isbombwalls {
// 			room.BroadcastMsg(usercmd.MsgTypeCmd_RetWallState, MapStateCmd(this))
// 		}

// 		this.IsBombPlayer(room)

// 		room.rangeBalls = room.rangeBalls[1:]
// 		this.player.bombLeft++
// 	}

// 	return true
// }

// //是否摧毁墙
// func (this *Bomb) IsBombWalls(room *Room) bool {
// 	power := this.power
// 	x := this.pos.X
// 	y := this.pos.Y
// 	flag := false
// 	// 上
// 	for i := y - 1; i >= y-power; i-- {
// 		if room.scene.gameMap.gamemap[x][i] == 0 {
// 			continue
// 		} else if room.scene.gameMap.gamemap[x][i] == 3 {
// 			// 摧毁墙
// 			room.scene.gameMap.gamemap[x][i] = 0
// 			flag = true
// 			this.boxDistroiedList = append(this.boxDistroiedList, &RetMapCellState{x, y, 0})
// 		}
// 	}
// 	// 下
// 	for i := y + 1; i <= y+power; i++ {
// 		if room.scene.gameMap.gamemap[x][i] == 0 {
// 			continue
// 		} else if room.scene.gameMap.gamemap[x][i] == 3 {
// 			// 摧毁墙
// 			room.scene.gameMap.gamemap[x][i] = 0
// 			flag = true
// 			this.boxDistroiedList = append(this.boxDistroiedList, &RetMapCellState{x, y, 0})

// 		}
// 	}
// 	// 左
// 	for i := x - 1; i >= x-power; i-- {
// 		if room.scene.gameMap.gamemap[i][y] == 0 {
// 			continue
// 		} else if room.scene.gameMap.gamemap[i][y] == 3 {
// 			// 摧毁墙
// 			room.scene.gameMap.gamemap[i][y] = 0
// 			flag = true
// 			this.boxDistroiedList = append(this.boxDistroiedList, &RetMapCellState{x, y, 0})

// 		}
// 	}

// 	// 右
// 	for i := x + 1; i <= x+power; i++ {
// 		if room.scene.gameMap.gamemap[i][y] == 0 {
// 			continue
// 		} else if room.scene.gameMap.gamemap[i][y] == 3 {
// 			// 摧毁墙
// 			room.scene.gameMap.gamemap[i][y] = 0
// 			flag = true
// 			this.boxDistroiedList = append(this.boxDistroiedList, &RetMapCellState{x, y, 0})

// 		}
// 	}

// 	return flag
// }

// //是否扎到人
// func (this *Bomb) IsBombPlayer(room *Room) bool {
// 	x := this.pos.X
// 	y := this.pos.Y
// 	flag := false
// 	for _, player := range room.Scene.players {
// 		if player.pos.x <= float64(x+this.power) && player.pos.x <= float64(x-this.power) && player.pos.y <= float64(y+this.power) && player.pos.y >= float64(y-this.power) {
// 			// 炸到了
// 			this.player.lifeState = 0

// 			// 广播该玩家被炸了
// 			room.BroadcastMsg(usercmd.MsgTypeCmd_PlayerState, PlayerStateCmd(player.id, int32(player.lifeState)))
// 			flag = true
// 		} else {
// 			continue
// 		}
// 	}

// 	return flag
// }

// func (this *Bomb) TimeAction() {

// }

// type BombMgr struct {
// 	room  *Room
// 	bombs []BaseBomb
// }

// func NewBombMgr(room *Room) *BombMgr {
// 	return &BombMgr{
// 		room: room,
// 	}
// }

// func (this *BombMgr) Add(bomb BaseBomb) bool {
// 	if !bomb.Init(this.room) {
// 		return false
// 	}
// 	this.bombs = append(this.bombs, bomb)
// 	return true
// }

// // 0.1秒执行一次
// func (this *BombMgr) ExecAction() {
// 	for i := 0; i < len(this.bombs); {
// 		bomb := this.bombs[i]
// 		if !bomb.Check(this.room) {
// 			i++
// 			continue
// 		}
// 		if !bomb.Exec(this.room) {
// 			i++
// 			continue
// 		}
// 		this.bombs = append(this.bombs[:i], this.bombs[i+1:]...)
// 	}
// }
