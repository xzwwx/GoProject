package main

import "sync"

type RoomMgr struct{
	runmutex sync.RWMutex
	runrooms map[uint32]*Room
}