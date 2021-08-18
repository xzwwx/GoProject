package common

//"encoding/binary"
//"glog"
//"net/http"
//"net/url"
//"strconv"
//"usercmd"

//Server Type
const (
	ServerTypeRoom  = 1 // Room Server
	ServerTypeTeam  = 2
	ServerTypeLogin = 3 // Login Server
)

// Cmd size
const (
	CmdHeaderSize = 2
	ServerCmdSize = 1
	ServerIdSize  = 4
	SubCmdSize    = 2
)

// Player room token data
type UserData struct {
	ServerId uint16 `redis:"ServerId"` // Vreify server Id
	Id       uint64 `redis:"Id"`       // userid
	Account  string `redis:"Account"`  // username

	RoomId   uint32 `redis:"RoomId"`
	RoomAddr string `redis:"RoomAddr"`
}

// Request message in game//////
type ReqMoveMsg struct {
	UserId    uint64
	Speed     uint32
	Direction uint32
}

type ReqLayBombMsg struct {
	UserId uint64
	X      uint32
	Y      uint32
}

// chi daoju
type ReqTriggerObjectMsg struct {
	UserId uint64
	ObjId  uint32 // daoju id
}

// bei zha
type ReqTriggerBombMsg struct {
	UserId uint64
	X      uint32
	Y      uint32
}

// kill player
type ReqKillMsg struct {
	UserId   uint64
	beKilled uint64
}

//////////////////

// Return from server

//Get Cmd
//func GetCmd(buf []byte)uint16{
//	if len(buf) <CmdHeaderSize{
//		return 0
//	}
//	return uint16(buf[0])|uint16(buf[1])<<8
//}

type RoomTokenInfo struct {
	UserId   uint64
	UserName string
	RoomId   uint32
}

type Vector2 struct {
	x float64
	y float64
}

// 格子
type GridPos struct {
	X uint32
	Y uint32
}

// 障碍物
type Obstacle struct {
	Id  uint32
	Pos GridPos
}

// 宝箱
type Box struct {
	Goods uint32 // 物体类型
	Id    uint32
	Pos   GridPos
}

// 坐标
type Position struct {
	X float64
	Y float64
}

type MoveWay byte

const (
	MoveWay_None = MoveWay(iota)
	MoveWay_Up
	MoveWay_Down
	MoveWay_Left
	MoveWay_Right
)
