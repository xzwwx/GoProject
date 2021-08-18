package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"glog"
)

type GridType uint32

const (
	GridType_Space    = GridType(iota) // 空地
	GridType_Obstacle                  // 不可摧毁的障碍物
	GridType_Box                       // 箱子（可摧毁）
	GridType_Bomb                      // 炸弹
)

type GameMap struct {
	Width    uint32
	Height   uint32
	MapArray [][]GridType
}

const (
	BGFilePath         = "./gamemap/BG.json"         // 空地，但是空地上有可能有障碍物
	BoundFilePath      = "./gamemap/Bound.json"      // 地图边界（暂时用不到）
	ForeGroundFilePath = "./gamemap/ForeGround.json" // 障碍物
)

const (
	OFFSET_X = 16
	OFFSET_Y = 8
)

func ReadFileAll(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func (this *GameMap) GetGridByPos(x, y uint32) GridType {
	return this.MapArray[x][y]
}

// 自定义地图信息
// 自定义地图信息
func (this *GameMap) CustomizeInit() bool {
	glog.Infoln("[游戏地图初始化] 初始化开始")
	this.Width, this.Height = 40, 40
	this.MapArray = make([][]GridType, this.Height)
	for i := 0; i < int(this.Height); i++ {
		this.MapArray[i] = make([]GridType, this.Width)
	}

	splitPos := func(posStr string) (int, int) {
		arr := strings.Split(posStr, ",")
		x, _ := strconv.Atoi(arr[0])
		y, _ := strconv.Atoi(arr[1])
		return x + OFFSET_X, y + OFFSET_Y
	}

	// 读取地图json文件，障碍物
	jsonBuf, err := ReadFileAll(ForeGroundFilePath)
	m := make(map[string]string)
	err = json.Unmarshal(jsonBuf, &m)
	if err != nil {
		glog.Errorln("[读取地图json文件错误] ", err)
		return false
	}
	for key, value := range m {
		_ = value
		x, y := splitPos(key)
		if x < 0 || x >= 40 || y < 0 || y >= 40 {
			// 坐标不合理
			continue
		}
		this.MapArray[x][y] = GridType_Obstacle
	}
	// 读取地图json文件，可移动的空地
	jsonBuf, err = ReadFileAll(BGFilePath)
	m = make(map[string]string)
	err = json.Unmarshal(jsonBuf, &m)
	if err != nil {
		glog.Errorln("[读取地图json文件错误] ", err)
		return false
	}
	for key, value := range m {
		_ = value
		x, y := splitPos(key)
		if x < 0 || x >= 40 || y < 0 || y >= 40 {
			// 坐标不合理
			continue
		}
		if this.MapArray[x][y] == GridType_Obstacle {
			continue
		}
		this.MapArray[x][y] = GridType_Space
	}
	glog.Infoln("[游戏地图初始化] 初始化完成")

	return true
}

// func (this *GameMap) CanPass(x, y int) bool {
// 	return this.MapArray[x][y] != GridType_Box &&
// 		this.MapArray[x][y] != GridType_Obstacle
// }

// func (this *GameMap) GetWidth() uint32 {
// 	return this.Width
// }

// func (this *GameMap) GetHeight() uint32 {
// 	return this.Height
// }

//package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"sync"
// 	"time"
// 	"usercmd"
// )

// type Map struct {
// 	gamemap   [][]int32
// 	mapPool   []map[uint32]interface{}
// 	needMutex bool
// 	mutexPool sync.RWMutex
// }

// type RetMapCellState struct {
// 	x     int32
// 	y     int32
// 	state int32
// }

// /*
// enum CellType{
// 	Space = 0;
// 	Wall = 1;
// 	Bomb = 2;
// 	Box  = 3;
// 	Object = 4;
//   }
// */
// func GenerateRandMap() (m usercmd.MapVector) {

// 	rand.Seed(time.Now().UnixNano())
// 	var sm []*usercmd.MapVector_Row

// 	for i := 0; i < 5; i++ {
// 		var row []usercmd.CellType
// 		for j := 0; j < 5; j++ {
// 			mtype := usercmd.CellType(rand.Intn(5))
// 			row = append(row, mtype)
// 			fmt.Println(mtype)
// 		}
// 		r := &usercmd.MapVector_Row{
// 			Y: row,
// 		}
// 		sm = append(sm, r)
// 		//sm = append(sm[:i], row)
// 	}
// 	m.X = sm
// 	return
// }

// //解析地图
// func DecodeMap(m usercmd.MapVector) [][]int32 {
// 	gamemap := make([][]int32, len(m.X))

// 	for i := range gamemap {
// 		gamemap[i] = make([]int32, len(m.X))
// 	}

// 	fmt.Println(len(m.X), " =======x==========")
// 	fmt.Println(len(m.X[0].Y), " =======y==========")

// 	for i := 0; i < len(m.X); i++ {
// 		for j := 0; j < len(m.X[0].Y); j++ {
// 			fmt.Println((m.X[i].Y[j]))
// 			ct := m.X[i].Y[j]
// 			var mt int32
// 			switch ct {
// 			case usercmd.CellType_Space:
// 				mt = 0
// 			case usercmd.CellType_Wall:
// 				mt = 1
// 			case usercmd.CellType_Bomb:
// 				mt = 2
// 			case usercmd.CellType_Box:
// 				mt = 3
// 			case usercmd.CellType_Object:
// 				mt = 4
// 			}
// 			gamemap[i][j] = mt
// 		}
// 	}
// 	return gamemap
// }

// //解析地图
// func DecodeMap0(m usercmd.MapVector) [][]usercmd.CellType {
// 	gamemap := make([][]usercmd.CellType, len(m.X))

// 	for i := range gamemap {
// 		gamemap[i] = make([]usercmd.CellType, len(m.X))
// 	}

// 	fmt.Println(len(m.X), " =======x==========")
// 	fmt.Println(len(m.X[0].Y), " =======y==========")

// 	for i := 0; i < len(m.X); i++ {
// 		for j := 0; j < len(m.X[0].Y); j++ {
// 			fmt.Println((m.X[i].Y[j]))

// 			gamemap[i][j] = m.X[i].Y[j]
// 		}
// 	}
// 	return gamemap
// }

// // 返回地图状态数组
// func MapStateCmd(bomb *Bomb) *usercmd.RetMapState {
// 	mapcell := &usercmd.RetMapState{}
// 	cc := mapcell.Cs
// 	for i := 0; i < len(bomb.boxDistroiedList); i++ {
// 		c := bomb.boxDistroiedList[i]
// 		cell := &usercmd.RetMapState_CellState{
// 			X:     c.x,
// 			Y:     c.y,
// 			State: usercmd.CellType(c.state),
// 		}
// 		cc = append(cc, cell)
// 	}

// 	return mapcell
// }

// // 道具
// type Objects struct {
// 	objType   int32 // 1: 加速 speed   2：威力 power  3：数量  Bombnum
// 	x         int32
// 	y         int32
// 	isExisted bool
// 	obj       [][]int32
// }

// //道具管理器
// type ObjMgr struct {
// 	room *Room
// 	objs []*Objects
// }

// // func main() {
// // 	m := GenerateRandMap()
// // 	gm := DecodeMap(m)
// // 	for i := 0; i < len(m.X); i++ {
// // 		for j := 0; j < len(m.X[0].Y); j++ {
// // 			fmt.Print(gm[i][j], " ")
// // 		}
// // 		fmt.Print("\n")
// // 	}
// // }
