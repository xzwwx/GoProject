package main

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"usercmd"
	"glog"
)

const (
	CmdHeaderSize 	= 2
	ServerCmdSize	= 1
)

// Get request parameters list
func GetPostValues(res http.ResponseWriter, req * http.Request)(cmd usercmd.MsgType, values url.Values, uid uint64, ver int, ok bool){
	qqttype, _ :=strconv.Atoi(req.Header.Get("qqt"))
	switch qqttype {
	case 2:

		body := make([]byte, req.ContentLength)
		blen, err := req.Body.Read(body)
		if err != nil || blen < 8 {
			glog.Error("[Login] Protocal too short. ", blen, ", ", body, ", ", err)
			return
		}
		uid = binary.LittleEndian.Uint64(body[blen - 8:])
		msg := &usercmd.ReqHttpArgData{}

		values = url.Values{}
		for _, arg := range msg.Args{
			values.Set(arg.Key, arg.Val)
		}

		return usercmd.MsgType(msg.Cmd), values, uid, qqttype, true
		//if err := msg.

	}

	return
}

// Get client ip:port
func getRemoteAddr(req *http.Request) string {
	if req == nil {
		return ""
	}
	if rip := req.Header.Get("X-QQDZ-IP"); rip != "" {
		return rip
	}
	return req.RemoteAddr
}

//Get IP
func GetIP(addr string) string {
	return strings.Split(addr,":")[0]
}

func sendCmd(res http.ResponseWriter, cmd usercmd.MsgType, msg proto.Message) bool {
	//json encoder
	retbuf := EncodeHttpCmd(cmd, msg)
	if retbuf == nil {
		return false
	}
	res.Write(retbuf)
	return true
}

// Encode Http Cmd
func EncodeHttpCmd(cmd usercmd.MsgType, msg proto.Message) []byte {
	// encode
	data, flag, err := EncodeCmd(uint16(cmd), msg);
	if err != nil {
		glog.ErrorDepth(1, "[CMD] Cmd encode failed. ", cmd)
		return nil
	}

	var retbuf []byte
	bsize := len(data)
	retbuf = append(retbuf, byte(bsize), byte(bsize >> 8), byte(bsize>>16), flag)
	retbuf = append(retbuf, data...)
	return retbuf
}

// encode
func EncodeCmd(cmd uint16, msg proto.Message)([]byte, byte, error){
	data , err := proto.Marshal(msg)
	if err != nil {
		glog.Error("[Protocal] Encode protobuf data failed. ", err)
		return nil, 0, err
	}
	var (
		mflag byte
		mbuff []byte
	)
	mflag = 0
	mbuff = data
	p := make([]byte, len(mbuff) + CmdHeaderSize)
	binary.LittleEndian.PutUint16(p[0:], cmd)
	copy(p[2:], mbuff)
	return p, mflag, nil

}