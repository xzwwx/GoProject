package main

import (
	"glog"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"usercmd"

	"github.com/golang/protobuf/proto"
)

// Get server time
func HandleTime(res http.ResponseWriter, req *http.Request) {
	timenow := time.Now().Unix()
	res.Write([]byte(strconv.FormatInt(timenow, 10)))
}

// Depose Login related message Format: /login?c=xxx
func HandleLoginMsg(res http.ResponseWriter, req *http.Request) {

}

// Depose game related message  Format: /game?c=xxx
func HandleGameMsg(res http.ResponseWriter, req *http.Request) {

	cmd, values, userid, _, ok := GetPostValues(res, req)
	if !ok {
		return
	}
	if userid == 0 {
		glog.Error("[Login] Session error.")
		return
	}
	switch cmd {
	case usercmd.MsgType_ReqIntoFRoom:
		{
			//srcip := GetIP(getRemoteAddr(req))

			// Into Room 	/game?c=ReqIntoFRoom&ver=2&model=0
			IntoRoom(res, req, values, userid)
		}
	}
}

// Depose messages between servers
func HandleServer(userid uint64, data []byte) {

}

// Into Room 	/game?c=ReqIntoFRoom&ver=2&model=0
func IntoRoom(res http.ResponseWriter, req *http.Request, values url.Values, userid uint64) {
	var (
	//ok 			bool
	//serverid	uint16
	//roomaddr	string
	//key			string
	//roomid		uint32
	//err 		error
	//onlineTime	int64
	//loginTime	int64
	//gServerId	uint32
	//roomEndTime uint32
	//seqIndex	uint32
	//
	//newsync		bool
	)
	//srcip := GetIP(getRemoteAddr(req))
	//location, cnet := 0, 0
	glog.Info("[Login] Received login request")

	// Free room
	//cityid := 0
	retCmd := &usercmd.RetIntoRoom{
		Err: proto.Uint32(388),
	}
	sendCmd(res, usercmd.MsgType_ReqIntoFRoom, retCmd)
	return

}
