package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
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

/*
enum CellType{
	Space = 0;
	Wall = 1;
	Bomb = 2;
	Box  = 3;
	Object = 4;
  }
*/
func GenerateRandMap() (m usercmd.MapVector) {

	rand.Seed(time.Now().UnixNano())
	var sm []*usercmd.MapVector_Row

	for i := 0; i < 5; i++ {
		var row []usercmd.CellType
		for j := 0; j < 5; j++ {
			mtype := usercmd.CellType(rand.Intn(5))
			row = append(row, mtype)
			fmt.Println(mtype)
		}
		r := &usercmd.MapVector_Row{
			Y: row,
		}
		sm = append(sm, r)
		//sm = append(sm[:i], row)
	}
	m.X = sm
	return
}

//解析地图
func DecodeMap(m usercmd.MapVector) [][]int32 {
	gamemap := make([][]int32, len(m.X))

	for i:= range gamemap {
		gamemap[i] = make([]int32, len(m.X))
	}

   fmt.Println(len(m.X)," =======x==========")
   fmt.Println(len(m.X[0].Y)," =======y==========")

   for i := 0; i < len(m.X); i++ {
	   for j := 0; j < len(m.X[0].Y); j++ {
		   fmt.Println((m.X[i].Y[j]))
			ct := m.X[i].Y[j]
			var mt int32
			switch ct {
			case usercmd.CellType_Space:
				mt = 0
			case usercmd.CellType_Wall:
				mt = 1
			case usercmd.CellType_Bomb:
				mt = 2
			case usercmd.CellType_Box:
				mt = 3
			case usercmd.CellType_Object:
				mt = 4		
			}
		   gamemap[i][j] = mt
	   }
   }
   return gamemap
}

//解析地图
func DecodeMap0(m usercmd.MapVector) [][]usercmd.CellType {
	 gamemap := make([][]usercmd.CellType, len(m.X))

	 for i:= range gamemap {
		 gamemap[i] = make([]usercmd.CellType, len(m.X))
	 }

	fmt.Println(len(m.X)," =======x==========")
	fmt.Println(len(m.X[0].Y)," =======y==========")

	for i := 0; i < len(m.X); i++ {
		for j := 0; j < len(m.X[0].Y); j++ {
			fmt.Println((m.X[i].Y[j]))

			gamemap[i][j] = m.X[i].Y[j]
		}
	}
	return gamemap
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
	obj       [][]int32
}

//道具管理器
type ObjMgr struct {
	room *Room
	objs []*Objects
}

// func main() {
// 	m := GenerateRandMap()
// 	gm := DecodeMap(m)
// 	for i := 0; i < len(m.X); i++ {
// 		for j := 0; j < len(m.X[0].Y); j++ {
// 			fmt.Print(gm[i][j], " ")
// 		}
// 		fmt.Print("\n")
// 	}
// }
