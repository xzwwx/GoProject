package main

import (
	"usercmd"
)

//player pool
type BallPool struct {
	moveFree          []*usercmd.MsgPlayerMove
	msgPlayer         []*usercmd.ScenePlayer
	msgScene          usercmd.MsgScene
	msgBomb           []*usercmd.MsgBomb
	RetUpdateSceneMsg *usercmd.RetUpdateSceneMsg
}
