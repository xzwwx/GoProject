package main

import (
	"time"
	"usercmd"
)

type BaseBomb interface {
	Init(room *Room) bool
	Check(room *Room) bool
	Exec(room *Room) bool
}

type Bomb struct {
	id               int64
	player           *ScenePlayer
	layTime          int64
	pos              VectorInt
	power            int32
	lastTime         int64 // 爆炸时间
	isdelete         bool  // zha le ma
	boxDistroiedList []*RetMapCellState
}

func (this *Bomb) Init(room *Room) bool {
	timenow := room.now.Unix()
	this.lastTime = timenow + time.Second.Milliseconds()*3 // 3秒后爆炸
	this.pos.X = int32(this.player.pos.x)
	this.pos.Y = int32(this.player.pos.y)
	this.power = this.player.self.power
	room.rangeBalls = append(room.rangeBalls, this)
	this.isdelete = false

	return true
}

func (this *Bomb) Check(room *Room) bool {
	timenow := time.Now().Unix()
	if timenow > this.lastTime { // 到爆炸时间了
		return true
	}
	return false

}

func (this *Bomb) Exec(room *Room) bool {
	// 先放的先炸
	bomb := room.rangeBalls[0]

	if bomb.layTime == this.layTime {

		isbombwalls := this.IsBombWalls(room)
		if isbombwalls {
			room.BroadcastMsg(usercmd.MsgTypeCmd_RetWallState, MapStateCmd(this))
		}

		this.IsBombPlayer(room)

		room.rangeBalls = room.rangeBalls[1:]
		this.player.bombLeft++
	}

	return true
}

//是否摧毁墙
func (this *Bomb) IsBombWalls(room *Room) bool {
	power := this.power
	x := this.pos.X
	y := this.pos.Y
	flag := false
	// 上
	for i := y - 1; i >= y-power; i-- {
		if room.scene.gameMap.gamemap[x][i] == 0 {
			continue
		} else if room.scene.gameMap.gamemap[x][i] == 3 {
			// 摧毁墙
			room.scene.gameMap.gamemap[x][i] = 0
			flag = true
			this.boxDistroiedList = append(this.boxDistroiedList, &RetMapCellState{x, y, 0})
		}
	}
	// 下
	for i := y + 1; i <= y+power; i++ {
		if room.scene.gameMap.gamemap[x][i] == 0 {
			continue
		} else if room.scene.gameMap.gamemap[x][i] == 3 {
			// 摧毁墙
			room.scene.gameMap.gamemap[x][i] = 0
			flag = true
			this.boxDistroiedList = append(this.boxDistroiedList, &RetMapCellState{x, y, 0})

		}
	}
	// 左
	for i := x - 1; i >= x-power; i-- {
		if room.scene.gameMap.gamemap[i][y] == 0 {
			continue
		} else if room.scene.gameMap.gamemap[i][y] == 3 {
			// 摧毁墙
			room.scene.gameMap.gamemap[i][y] = 0
			flag = true
			this.boxDistroiedList = append(this.boxDistroiedList, &RetMapCellState{x, y, 0})

		}
	}

	// 右
	for i := x + 1; i <= x+power; i++ {
		if room.scene.gameMap.gamemap[i][y] == 0 {
			continue
		} else if room.scene.gameMap.gamemap[i][y] == 3 {
			// 摧毁墙
			room.scene.gameMap.gamemap[i][y] = 0
			flag = true
			this.boxDistroiedList = append(this.boxDistroiedList, &RetMapCellState{x, y, 0})

		}
	}

	return flag
}

//是否扎到人
func (this *Bomb) IsBombPlayer(room *Room) bool {
	x := this.pos.X
	y := this.pos.Y
	flag := false
	for _, player := range room.Scene.players {
		if player.pos.x <= float64(x+this.power) && player.pos.x <= float64(x-this.power) && player.pos.y <= float64(y+this.power) && player.pos.y >= float64(y-this.power) {
			// 炸到了
			this.player.lifeState = 0

			// 广播该玩家被炸了
			room.BroadcastMsg(usercmd.MsgTypeCmd_PlayerState, PlayerStateCmd(player.id, int32(player.lifeState)))
			flag = true
		} else {
			continue
		}
	}

	return flag
}

func (this *Bomb) TimeAction() {

}

type BombMgr struct {
	room  *Room
	bombs []BaseBomb
}

func NewBombMgr(room *Room) *BombMgr {
	return &BombMgr{
		room: room,
	}
}

func (this *BombMgr) Add(bomb BaseBomb) bool {
	if !bomb.Init(this.room) {
		return false
	}
	this.bombs = append(this.bombs, bomb)
	return true
}

// 0.1秒执行一次
func (this *BombMgr) ExecAction() {
	for i := 0; i < len(this.bombs); {
		bomb := this.bombs[i]
		if !bomb.Check(this.room) {
			i++
			continue
		}
		if !bomb.Exec(this.room) {
			i++
			continue
		}
		this.bombs = append(this.bombs[:i], this.bombs[i+1:]...)
	}
}
