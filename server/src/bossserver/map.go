package main

import (
	"sync"
	"usercmd"
)

type Map struct {
	gamemap   [][]int32
	mapPool   []map[uint32]interface{}
	needMutex bool
	mutexPool sync.RWMutex
}

type RetMapCellState struct {
	x     int32
	y     int32
	state int32
}

func GenerateRandMap() *map[uint32]*Obstacle {
	obstacle := make(map[uint32]*Obstacle)
	for i := 0; i < 30; i++ {
		o := &Obstacle{
			Id:    uint32(i),
			pos:   VectorInt{0, 0},
			OType: 1,
		}
		obstacle[uint32(i)] = o
	}
	return &obstacle

}

// 返回地图状态数组
func MapStateCmd(bomb *Bomb) *usercmd.RetMapState {
	mapcell := &usercmd.RetMapState{}
	cc := mapcell.Cs
	for i := 0; i < len(bomb.boxDistroiedList); i++ {
		c := bomb.boxDistroiedList[i]
		cell := &usercmd.RetMapState_CellState{
			X:     c.x,
			Y:     c.y,
			State: usercmd.CellType(c.state),
		}
		cc = append(cc, cell)
	}

	return mapcell
}

// 道具
type Objects struct {
	objType   int32 // 1: 加速 speed   2：威力 power  3：数量  Bombnum
	x         int32
	y         int32
	isExisted bool
	obj 	[][]int32
}

//道具管理器
type ObjMgr struct {
	room *Room
	objs []*Objects
}
